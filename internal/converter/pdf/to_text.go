package pdf

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts a PDF document to the specified output format
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	// Route to appropriate converter based on output format
	switch opts.OutputFormat {
	case internal.FormatTXT:
		return c.toText(ctx, input, output, opts)
	case internal.FormatPNG, internal.FormatJPG, internal.FormatJPEG:
		return c.toImage(ctx, input, output, opts)
	default:
		return fmt.Errorf("%w: PDF to %s", internal.ErrUnsupportedConversion, opts.OutputFormat)
	}
}

// toText extracts text from PDF and writes it to output file
func (c *Converter) toText(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	log.Info().
		Str("input", input).
		Str("output", output).
		Str("engine", c.engine).
		Msg("Converting PDF to text")

	// Create PDF configuration
	conf := model.NewDefaultConfiguration()

	// Create temporary directory for extracted pages
	tmpDir, err := os.MkdirTemp("", "yakateka-pdf-*")
	if err != nil {
		return fmt.Errorf("%w: failed to create temp directory: %v", internal.ErrConversionFailed, err)
	}
	defer os.RemoveAll(tmpDir)

	// Extract text from all pages using pdfcpu API
	// ExtractContentFile writes separate text files per page to a directory
	err = api.ExtractContentFile(input, tmpDir, nil, conf)
	if err != nil {
		log.Error().Err(err).Str("input", input).Msg("Failed to extract text from PDF")
		return fmt.Errorf("%w: %v", internal.ErrConversionFailed, err)
	}

	// Combine all extracted text files into a single output file
	err = combineTextFiles(tmpDir, output)
	if err != nil {
		log.Error().Err(err).Str("output", output).Msg("Failed to combine text files")
		return fmt.Errorf("%w: %v", internal.ErrConversionFailed, err)
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
		Msg("Successfully converted PDF to text")

	return nil
}

// combineTextFiles combines all .txt files in a directory into a single output file
func combineTextFiles(dir, output string) error {
	// Find all .txt files in the directory
	files, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no text files found in directory")
	}

	// Sort files by name to maintain page order
	sort.Strings(files)

	// Create output file
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Combine all files
	for i, file := range files {
		// Read input file
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", file, err)
		}

		// Write to output file
		if _, err := outFile.Write(content); err != nil {
			return fmt.Errorf("failed to write to output: %v", err)
		}

		// Add page separator (except for last page)
		if i < len(files)-1 {
			if _, err := io.WriteString(outFile, "\n\n--- Page Break ---\n\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

// toImage converts PDF pages to images (placeholder for future implementation)
func (c *Converter) toImage(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	return fmt.Errorf("%w: PDF to image conversion not yet implemented", internal.ErrUnsupportedConversion)
}
