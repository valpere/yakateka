package calibre

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts ebooks using Calibre's ebook-convert
// Usage: ebook-convert input.mobi output.epub [options]
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
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

	log.Info().
		Str("input", absInput).
		Str("output", absOutput).
		Str("from", string(opts.InputFormat)).
		Str("to", string(opts.OutputFormat)).
		Str("ebook-convert", c.ebookConvertPath).
		Msg("Converting ebook with Calibre")

	// Build ebook-convert command
	// Usage: ebook-convert input.mobi output.epub [options]
	args := []string{
		absInput,
		absOutput,
	}

	// Add optional arguments based on options
	if opts.Quality != "" {
		// Map quality to Calibre options
		switch opts.Quality {
		case "high":
			args = append(args, "--pretty-print")
		}
	}

	// Execute ebook-convert
	cmd := exec.CommandContext(ctx, c.ebookConvertPath, args...)

	// Capture output for logging
	outputBytes, err := cmd.CombinedOutput()
	outputStr := string(outputBytes)

	if err != nil {
		log.Error().
			Err(err).
			Str("input", absInput).
			Str("output", outputStr).
			Msg("Calibre ebook conversion failed")
		return fmt.Errorf("%w: Calibre ebook-convert failed: %v - %s",
			internal.ErrConversionFailed, err, outputStr)
	}

	// Verify output file was created
	if _, err := os.Stat(absOutput); os.IsNotExist(err) {
		log.Error().Str("output", absOutput).Msg("Output file was not created by Calibre")
		return fmt.Errorf("%w: output file not created", internal.ErrConversionFailed)
	}

	// Get output file size for logging
	stat, _ := os.Stat(absOutput)
	var fileSize int64
	if stat != nil {
		fileSize = stat.Size()
	}

	log.Info().
		Str("output", absOutput).
		Int64("bytes", fileSize).
		Str("conversion", fmt.Sprintf("%s → %s", opts.InputFormat, opts.OutputFormat)).
		Msg("Successfully converted ebook with Calibre")

	return nil
}

// CheckAvailability checks if Calibre ebook-convert is available on the system
func (c *Converter) CheckAvailability() error {
	if _, err := exec.LookPath(c.ebookConvertPath); err != nil {
		return fmt.Errorf("Calibre ebook-convert not available: %w", err)
	}
	return nil
}

// GetVersion returns the Calibre version
func (c *Converter) GetVersion() (string, error) {
	cmd := exec.Command(c.ebookConvertPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", nil
	}
	return string(output), nil
}
