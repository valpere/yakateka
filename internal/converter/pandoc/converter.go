package pandoc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts a document using Pandoc
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	log.Info().
		Str("input", input).
		Str("output", output).
		Str("from", string(opts.InputFormat)).
		Str("to", string(opts.OutputFormat)).
		Str("pandoc", c.pandocPath).
		Msg("Converting document with Pandoc")

	// Build pandoc command
	args := c.buildArgs(input, output, opts)

	// Execute pandoc
	cmd := exec.CommandContext(ctx, c.pandocPath, args...)

	// Capture output for logging
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("input", input).
			Str("output", string(output_bytes)).
			Msg("Pandoc conversion failed")
		return fmt.Errorf("%w: pandoc failed: %v - %s", internal.ErrConversionFailed, err, string(output_bytes))
	}

	// Verify output file was created
	if _, err := os.Stat(output); os.IsNotExist(err) {
		log.Error().Str("output", output).Msg("Output file was not created by Pandoc")
		return fmt.Errorf("%w: output file not created", internal.ErrConversionFailed)
	}

	// Get output file size for logging
	stat, _ := os.Stat(output)
	var fileSize int64
	if stat != nil {
		fileSize = stat.Size()
	}

	log.Info().
		Str("output", output).
		Int64("bytes", fileSize).
		Msg("Successfully converted document with Pandoc")

	return nil
}

// buildArgs constructs pandoc command-line arguments
func (c *Converter) buildArgs(input, output string, opts internal.ConversionOptions) []string {
	args := []string{
		input,
		"-o", output,
	}

	// Add format specifications if provided
	if opts.InputFormat != "" {
		args = append(args, "--from", toPandocFormat(opts.InputFormat))
	}
	if opts.OutputFormat != "" {
		args = append(args, "--to", toPandocFormat(opts.OutputFormat))
	}

	// Add quality-based options
	if opts.OutputFormat == internal.FormatPDF {
		// PDF-specific options based on quality
		switch opts.Quality {
		case "high":
			args = append(args, "--pdf-engine=xelatex")
		case "medium", "low":
			args = append(args, "--pdf-engine=pdflatex")
		}
	}

	// Add any extra arguments from config
	if len(c.extraArgs) > 0 {
		args = append(args, c.extraArgs...)
	}

	// Add extra arguments from options
	for key, value := range opts.Extra {
		args = append(args, fmt.Sprintf("--%s=%s", key, value))
	}

	return args
}

// CheckAvailability checks if Pandoc is available on the system
func (c *Converter) CheckAvailability() error {
	cmd := exec.Command(c.pandocPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pandoc not available at %s: %w", c.pandocPath, err)
	}
	return nil
}

// GetVersion returns the Pandoc version
func (c *Converter) GetVersion() (string, error) {
	cmd := exec.Command(c.pandocPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get pandoc version: %w", err)
	}

	// Parse version from output (first line usually contains version)
	lines := string(output)
	if len(lines) > 0 {
		return filepath.Base(lines), nil
	}

	return "unknown", nil
}

// toPandocFormat converts internal format names to Pandoc format names
func toPandocFormat(format internal.DocumentFormat) string {
	// Map internal format names to Pandoc format names
	switch format {
	case internal.FormatMD:
		return "markdown"
	case internal.FormatTXT:
		// Note: Pandoc can write "plain" text but cannot read it as input
		// When TXT is input format, we should not use Pandoc
		return "plain"
	case internal.FormatDOCX:
		return "docx"
	case internal.FormatODT:
		return "odt"
	case internal.FormatRTF:
		return "rtf"
	case internal.FormatHTML:
		return "html"
	case internal.FormatPDF:
		return "pdf"
	case internal.FormatEPUB:
		return "epub"
	case internal.FormatFB2:
		return "fb2"
	case internal.FormatJSON:
		return "json"
	case internal.FormatCSV:
		return "csv"
	case internal.FormatLaTeX:
		return "latex"
	case internal.FormatRST:
		return "rst"
	default:
		// Fall back to the internal format name (works for most formats)
		return string(format)
	}
}
