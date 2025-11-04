package djvu

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts a DjVu document to text or PostScript
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	// Route to appropriate converter based on output format
	switch opts.OutputFormat {
	case internal.FormatTXT, "":
		return c.convertToText(ctx, input, output)
	case internal.FormatPS:
		return c.convertToPS(ctx, input, output)
	default:
		return fmt.Errorf("%w: DjVu converter only supports TXT and PS output, requested %s",
			internal.ErrUnsupportedFormat, opts.OutputFormat)
	}
}

// convertToText extracts text from DjVu using djvutxt
func (c *Converter) convertToText(ctx context.Context, input, output string) error {

	log.Info().
		Str("input", input).
		Str("output", output).
		Str("djvutxt", c.djvutxtPath).
		Msg("Extracting text from DjVu document")

	// Build djvutxt command
	// Usage: djvutxt [options] <djvufile> [<outputfile>]
	args := []string{input, output}

	// Execute djvutxt
	cmd := exec.CommandContext(ctx, c.djvutxtPath, args...)

	// Capture output for logging
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("input", input).
			Str("output", string(outputBytes)).
			Msg("DjVu text extraction failed")
		return fmt.Errorf("%w: djvutxt failed: %v - %s",
			internal.ErrConversionFailed, err, string(outputBytes))
	}

	// Verify output file was created
	if _, err := os.Stat(output); os.IsNotExist(err) {
		log.Error().Str("output", output).Msg("Output file was not created by djvutxt")
		return fmt.Errorf("%w: output file not created", internal.ErrConversionFailed)
	}

	// Get output file size for logging
	stat, _ := os.Stat(output)
	var fileSize int64
	if stat != nil {
		fileSize = stat.Size()
	}

	// Warn if extraction is empty (no text layer)
	if fileSize == 0 {
		log.Warn().
			Str("input", input).
			Str("output", output).
			Msg("DjVu text extraction produced empty output - file may not have embedded text layer (requires OCR)")
		// Still return success, as djvutxt completed without error
		// Future: Could trigger OCR pipeline here
	} else {
		log.Info().
			Str("output", output).
			Int64("bytes", fileSize).
			Msg("Successfully extracted text from DjVu document")
	}

	return nil
}

// convertToPS converts DjVu to PostScript using djvups
func (c *Converter) convertToPS(ctx context.Context, input, output string) error {
	log.Info().
		Str("input", input).
		Str("output", output).
		Str("djvups", c.djvupsPath).
		Msg("Converting DjVu to PostScript")

	// Build djvups command
	// Usage: djvups [options] input.djvu output.ps
	args := []string{input, output}

	// Execute djvups
	cmd := exec.CommandContext(ctx, c.djvupsPath, args...)

	// Capture output for logging
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("input", input).
			Str("output", string(outputBytes)).
			Msg("DjVu to PS conversion failed")
		return fmt.Errorf("%w: djvups failed: %v - %s",
			internal.ErrConversionFailed, err, string(outputBytes))
	}

	// Verify output file was created
	if _, err := os.Stat(output); os.IsNotExist(err) {
		log.Error().Str("output", output).Msg("Output file was not created by djvups")
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
		Msg("Successfully converted DjVu to PostScript")

	return nil
}

// CheckAvailability checks if djvutxt is available on the system
func (c *Converter) CheckAvailability() error {
	// djvutxt exits with error when called without arguments, but prints usage
	// We just check if the binary exists in PATH
	if _, err := exec.LookPath(c.djvutxtPath); err != nil {
		return fmt.Errorf("djvutxt not available: %w", err)
	}
	return nil
}

// GetVersion returns the DjVuLibre version
func (c *Converter) GetVersion() (string, error) {
	cmd := exec.Command(c.djvutxtPath)
	output, _ := cmd.CombinedOutput()

	// Parse version from output (first line contains version)
	// Example: "DDJVU --- DjVuLibre-3.5.28"
	lines := string(output)
	if len(lines) > 0 {
		// Extract version from first line
		for i := 0; i < len(lines) && i < 100; i++ {
			if lines[i] == '\n' {
				return lines[:i], nil
			}
		}
		if len(lines) < 100 {
			return lines, nil
		}
		return lines[:100], nil
	}

	return "unknown", nil
}
