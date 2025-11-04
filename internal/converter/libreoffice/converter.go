package libreoffice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts documents using LibreOffice
// Usage: soffice --headless --convert-to <format> --outdir <dir> <input>
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	// Get output format and filter
	outputFormat, filterName := c.getOutputFormatAndFilter(opts.OutputFormat)

	// Get absolute paths
	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to get absolute input path: %w", err)
	}

	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}

	// Get output directory
	outDir := filepath.Dir(absOutput)

	log.Info().
		Str("input", absInput).
		Str("output", absOutput).
		Str("format", outputFormat).
		Str("filter", filterName).
		Str("soffice", c.sofficePath).
		Msg("Converting document with LibreOffice")

	// Build LibreOffice command
	// soffice --headless --convert-to <format>:<filter> --outdir <dir> <input>
	args := []string{
		"--headless",        // Run without GUI
		"--convert-to",      // Conversion mode
		outputFormat,        // Output format with optional filter
		"--outdir", outDir,  // Output directory
		absInput,            // Input file
	}

	// Execute LibreOffice conversion
	cmd := exec.CommandContext(ctx, c.sofficePath, args...)

	// Capture output for logging
	outputBytes, err := cmd.CombinedOutput()
	outputStr := string(outputBytes)

	// Filter out Python SyntaxWarning spam from LibreOffice extensions
	lines := strings.Split(outputStr, "\n")
	var filteredOutput []string
	for _, line := range lines {
		if !strings.Contains(line, "SyntaxWarning") &&
			!strings.Contains(line, "invalid escape sequence") &&
			!strings.Contains(line, "lightproof") {
			filteredOutput = append(filteredOutput, line)
		}
	}
	cleanOutput := strings.Join(filteredOutput, "\n")

	if err != nil {
		log.Error().
			Err(err).
			Str("input", absInput).
			Str("output", cleanOutput).
			Msg("LibreOffice conversion failed")
		return fmt.Errorf("%w: LibreOffice conversion failed: %v - %s",
			internal.ErrConversionFailed, err, cleanOutput)
	}

	// LibreOffice creates the file with the basename of input + new extension
	// We need to rename it to the desired output path
	baseName := strings.TrimSuffix(filepath.Base(absInput), filepath.Ext(absInput))
	expectedOutput := filepath.Join(outDir, baseName+"."+string(opts.OutputFormat))

	// Check if LibreOffice created the expected output
	if _, err := os.Stat(expectedOutput); err == nil {
		// Rename to desired output path if different
		if expectedOutput != absOutput {
			if err := os.Rename(expectedOutput, absOutput); err != nil {
				log.Error().
					Str("expected", expectedOutput).
					Str("desired", absOutput).
					Err(err).
					Msg("Failed to rename output file")
				return fmt.Errorf("failed to rename output: %w", err)
			}
		}
	} else {
		// Check if output was created at desired path directly
		if _, err := os.Stat(absOutput); os.IsNotExist(err) {
			log.Error().Str("output", absOutput).Msg("Output file was not created by LibreOffice")
			return fmt.Errorf("%w: output file not created", internal.ErrConversionFailed)
		}
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
		Msg("Successfully converted document with LibreOffice")

	return nil
}

// getOutputFormatAndFilter returns the LibreOffice output format string
// Format can be: "pdf" or "pdf:writer_pdf_Export" (format:filter)
func (c *Converter) getOutputFormatAndFilter(format internal.DocumentFormat) (string, string) {
	switch format {
	case internal.FormatPDF:
		return "pdf:writer_pdf_Export", "writer_pdf_Export"
	case internal.FormatHTML:
		return "html", "" // LibreOffice auto-selects appropriate HTML filter
	case internal.FormatTXT:
		return "txt:Text (encoded):UTF8", "Text (encoded):UTF8"
	case internal.FormatDOCX:
		return "docx", ""
	case internal.FormatODT:
		return "odt", ""
	case internal.FormatRTF:
		return "rtf", ""
	case internal.FormatMD:
		// LibreOffice doesn't have native Markdown export, use text
		log.Warn().Msg("LibreOffice doesn't support native Markdown export, using plain text")
		return "txt:Text (encoded):UTF8", "Text (encoded):UTF8"
	default:
		return string(format), ""
	}
}

// CheckAvailability checks if LibreOffice is available on the system
func (c *Converter) CheckAvailability() error {
	if _, err := exec.LookPath(c.sofficePath); err != nil {
		return fmt.Errorf("LibreOffice not available: %w", err)
	}
	return nil
}

// GetVersion returns the LibreOffice version
func (c *Converter) GetVersion() (string, error) {
	cmd := exec.Command(c.sofficePath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", nil
	}
	return strings.TrimSpace(string(output)), nil
}
