package postscript

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts PostScript to PDF using ps2pdf (Ghostscript)
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	log.Info().
		Str("input", input).
		Str("output", output).
		Str("ps2pdf", c.ps2pdfPath).
		Msg("Converting PostScript to PDF")

	// Build ps2pdf command
	// Usage: ps2pdf input.ps output.pdf
	args := []string{input, output}

	// Execute ps2pdf
	cmd := exec.CommandContext(ctx, c.ps2pdfPath, args...)

	// Capture output for logging
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("input", input).
			Str("output", string(outputBytes)).
			Msg("PostScript to PDF conversion failed")
		return fmt.Errorf("%w: ps2pdf failed: %v - %s",
			internal.ErrConversionFailed, err, string(outputBytes))
	}

	// Verify output file was created
	if _, err := os.Stat(output); os.IsNotExist(err) {
		log.Error().Str("output", output).Msg("Output file was not created by ps2pdf")
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
		Msg("Successfully converted PostScript to PDF")

	return nil
}

// CheckAvailability checks if ps2pdf is available on the system
func (c *Converter) CheckAvailability() error {
	if _, err := exec.LookPath(c.ps2pdfPath); err != nil {
		return fmt.Errorf("ps2pdf not available: %w", err)
	}
	return nil
}

// GetVersion returns the Ghostscript version
func (c *Converter) GetVersion() (string, error) {
	cmd := exec.Command("gs", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", nil
	}
	return string(output), nil
}
