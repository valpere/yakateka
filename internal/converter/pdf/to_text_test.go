package pdf

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/valpere/yakateka/internal"
)

func TestNewConverter(t *testing.T) {
	// Test with default engine
	c := NewConverter("")
	if c.engine != "pdfcpu" {
		t.Errorf("Expected default engine 'pdfcpu', got '%s'", c.engine)
	}

	// Test with specified engine
	c = NewConverter("unipdf")
	if c.engine != "unipdf" {
		t.Errorf("Expected engine 'unipdf', got '%s'", c.engine)
	}
}

func TestSupportedFormats(t *testing.T) {
	c := NewConverter("")

	// Test input formats
	inputFormats := c.SupportedInputFormats()
	if len(inputFormats) != 1 {
		t.Errorf("Expected 1 input format, got %d", len(inputFormats))
	}
	if inputFormats[0] != internal.FormatPDF {
		t.Errorf("Expected PDF input format, got %s", inputFormats[0])
	}

	// Test output formats
	outputFormats := c.SupportedOutputFormats()
	if len(outputFormats) != 3 {
		t.Errorf("Expected 3 output formats, got %d", len(outputFormats))
	}

	expectedFormats := map[internal.DocumentFormat]bool{
		internal.FormatTXT: true,
		internal.FormatPNG: true,
		internal.FormatJPG: true,
	}

	for _, format := range outputFormats {
		if !expectedFormats[format] {
			t.Errorf("Unexpected output format: %s", format)
		}
	}
}

func TestConvertInvalidInput(t *testing.T) {
	c := NewConverter("")
	ctx := context.Background()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatPDF,
		OutputFormat: internal.FormatTXT,
	}

	// Test with non-existent input file
	err := c.Convert(ctx, "nonexistent.pdf", "output.txt", opts)
	if err == nil {
		t.Error("Expected error for non-existent input file")
	}
}

func TestConvertUnsupportedFormat(t *testing.T) {
	c := NewConverter("")
	ctx := context.Background()

	// Create a temporary input file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.pdf")
	if err := os.WriteFile(inputFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatPDF,
		OutputFormat: internal.FormatDOCX, // Unsupported
	}

	// Test with unsupported output format
	err := c.Convert(ctx, inputFile, "output.docx", opts)
	if err == nil {
		t.Error("Expected error for unsupported output format")
	}
}

func TestToImageNotImplemented(t *testing.T) {
	c := NewConverter("")
	ctx := context.Background()

	// Create a temporary input file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.pdf")
	if err := os.WriteFile(inputFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatPDF,
		OutputFormat: internal.FormatPNG,
	}

	// Test that image conversion is not yet implemented
	err := c.Convert(ctx, inputFile, "output.png", opts)
	if err == nil {
		t.Error("Expected error for unimplemented image conversion")
	}
}
