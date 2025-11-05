package generic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
	"github.com/valpere/yakateka/internal/converter/config"
)

// Converter is a generic converter that executes commands based on configuration
type Converter struct {
	name     string
	config   config.ToolConfig
	profiles map[string]config.ProfileConfig
}

// NewConverter creates a new generic converter from configuration
func NewConverter(name string, cfg config.ToolConfig, profiles map[string]config.ProfileConfig) *Converter {
	return &Converter{
		name:     name,
		config:   cfg,
		profiles: profiles,
	}
}

// SupportedInputFormats returns formats this converter can read
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	formats := make([]internal.DocumentFormat, len(c.config.Formats.Input))
	for i, f := range c.config.Formats.Input {
		formats[i] = internal.DocumentFormat(f)
	}
	return formats
}

// SupportedOutputFormats returns formats this converter can write
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	formats := make([]internal.DocumentFormat, len(c.config.Formats.Output))
	for i, f := range c.config.Formats.Output {
		formats[i] = internal.DocumentFormat(f)
	}
	return formats
}

// Convert performs document conversion using configured command
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	// Check if conversion is supported
	if !c.config.SupportsConversion(string(opts.InputFormat), string(opts.OutputFormat)) {
		return fmt.Errorf("%w: %s does not support %s -> %s",
			internal.ErrUnsupportedConversion, c.name, opts.InputFormat, opts.OutputFormat)
	}

	// Get absolute paths
	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to get absolute input path: %w", err)
	}

	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}

	// Build command
	cmdStr, err := c.buildCommand(absInput, absOutput, opts)
	if err != nil {
		return fmt.Errorf("failed to build command: %w", err)
	}

	log.Info().
		Str("converter", c.name).
		Str("input", absInput).
		Str("output", absOutput).
		Str("from", string(opts.InputFormat)).
		Str("to", string(opts.OutputFormat)).
		Str("command", cmdStr).
		Msg("Converting document with generic converter")

	// Execute command
	outputBytes, err := c.executeCommand(ctx, cmdStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("converter", c.name).
			Str("output", string(outputBytes)).
			Msg("Conversion failed")
		return fmt.Errorf("%w: %s conversion failed: %v - %s",
			internal.ErrConversionFailed, c.name, err, string(outputBytes))
	}

	// Post-process if needed
	if err := c.postProcess(absInput, absOutput); err != nil {
		return fmt.Errorf("post-processing failed: %w", err)
	}

	// Verify output file was created
	if _, err := os.Stat(absOutput); os.IsNotExist(err) {
		log.Error().Str("output", absOutput).Msg("Output file was not created")
		return fmt.Errorf("%w: output file not created", internal.ErrConversionFailed)
	}

	// Get output file size for logging
	stat, _ := os.Stat(absOutput)
	var fileSize int64
	if stat != nil {
		fileSize = stat.Size()
	}

	log.Info().
		Str("converter", c.name).
		Str("output", absOutput).
		Int64("bytes", fileSize).
		Str("conversion", fmt.Sprintf("%s → %s", opts.InputFormat, opts.OutputFormat)).
		Msg("Successfully converted document")

	return nil
}

// buildCommand constructs the command string from template
func (c *Converter) buildCommand(input, output string, opts internal.ConversionOptions) (string, error) {
	template := c.config.GetCommandTemplate(c.profiles)
	if template == "" {
		return "", fmt.Errorf("no command template defined for %s", c.name)
	}

	// Get conversion override
	override := c.config.GetConversionOverride(string(opts.InputFormat), string(opts.OutputFormat))

	// Build replacements map
	replacements := map[string]string{
		"{binary}":        c.config.Binary,
		"{input}":         input,
		"{output}":        output,
		"{outdir}":        filepath.Dir(output),
		"{input_format}":  c.config.MapFormat(string(opts.InputFormat)),
		"{output_format}": c.config.MapFormat(string(opts.OutputFormat)),
		"{extra_args}":    "",
	}

	// Add override-specific replacements
	if override != nil {
		if override.ExtraArgs != "" {
			replacements["{extra_args}"] = override.ExtraArgs
		}
		if override.OutputFormat != "" {
			replacements["{output_format}"] = override.OutputFormat
			replacements["{format}"] = override.OutputFormat
		}

		// Add quality flags if specified
		if opts.Quality != "" {
			if qualityFlags, ok := override.Quality[opts.Quality]; ok {
				if replacements["{extra_args}"] != "" {
					replacements["{extra_args}"] += " " + qualityFlags
				} else {
					replacements["{extra_args}"] = qualityFlags
				}
			}
		}
	}

	// Also support {format} as alias for {output_format}
	if _, ok := replacements["{format}"]; !ok {
		replacements["{format}"] = replacements["{output_format}"]
	}

	// Replace all placeholders
	cmdStr := template
	for placeholder, value := range replacements {
		cmdStr = strings.ReplaceAll(cmdStr, placeholder, value)
	}

	// Clean up extra spaces
	cmdStr = strings.Join(strings.Fields(cmdStr), " ")

	return cmdStr, nil
}

// executeCommand executes the command string
// Uses shlex to properly handle quoted arguments like --option="value with spaces"
func (c *Converter) executeCommand(ctx context.Context, cmdStr string) ([]byte, error) {
	// Split command into parts using shell word splitting
	// This properly handles quoted strings and escaped characters
	parts, err := shlex.Split(cmdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	return cmd.CombinedOutput()
}

// postProcess handles post-conversion processing (e.g., file renaming)
func (c *Converter) postProcess(input, output string) error {
	postProcess := c.config.GetPostProcess(c.profiles)

	switch postProcess {
	case "rename_from_basename":
		// LibreOffice-style: it creates basename.ext instead of desired output path
		baseName := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		outputExt := filepath.Ext(output)
		expectedOutput := filepath.Join(filepath.Dir(output), baseName+outputExt)

		// Check if file was created with basename
		if _, err := os.Stat(expectedOutput); err == nil {
			// Rename to desired output path if different
			if expectedOutput != output {
				if err := os.Rename(expectedOutput, output); err != nil {
					return fmt.Errorf("failed to rename output: %w", err)
				}
			}
		}
	}

	return nil
}
