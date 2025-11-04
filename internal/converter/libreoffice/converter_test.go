package libreoffice

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/valpere/yakateka/internal"
)

func TestNewConverter(t *testing.T) {
	t.Run("custom_path", func(t *testing.T) {
		c := NewConverter("/custom/path/soffice")
		if c.sofficePath != "/custom/path/soffice" {
			t.Errorf("Expected /custom/path/soffice, got %s", c.sofficePath)
		}
	})

	t.Run("empty_path_defaults_to_PATH", func(t *testing.T) {
		c := NewConverter("")
		if c.sofficePath != "soffice" {
			t.Errorf("Expected soffice, got %s", c.sofficePath)
		}
	})
}

func TestSupportedFormats(t *testing.T) {
	c := NewConverter("")

	inputFormats := c.SupportedInputFormats()
	if len(inputFormats) == 0 {
		t.Error("Expected non-zero input formats")
	}

	// Check for key formats
	hasPS := false
	hasDOCX := false
	for _, f := range inputFormats {
		if f == internal.FormatPS {
			hasPS = true
		}
		if f == internal.FormatDOCX {
			hasDOCX = true
		}
	}
	if !hasPS {
		t.Error("Expected PS in input formats")
	}
	if !hasDOCX {
		t.Error("Expected DOCX in input formats")
	}

	outputFormats := c.SupportedOutputFormats()
	if len(outputFormats) == 0 {
		t.Error("Expected non-zero output formats")
	}

	// Check for key formats
	hasHTML := false
	hasPDF := false
	hasMD := false
	for _, f := range outputFormats {
		if f == internal.FormatHTML {
			hasHTML = true
		}
		if f == internal.FormatPDF {
			hasPDF = true
		}
		if f == internal.FormatMD {
			hasMD = true
		}
	}
	if !hasHTML {
		t.Error("Expected HTML in output formats")
	}
	if !hasPDF {
		t.Error("Expected PDF in output formats")
	}
	if hasMD {
		t.Error("MD should NOT be in output formats (use HTML→MD via Pandoc)")
	}
}

func TestGetOutputFormatAndFilter(t *testing.T) {
	c := NewConverter("")

	tests := []struct {
		name           string
		format         internal.DocumentFormat
		wantFormat     string
		wantFilterPart string
	}{
		{"PDF", internal.FormatPDF, "pdf:writer_pdf_Export", "writer_pdf_Export"},
		{"HTML", internal.FormatHTML, "html", ""},
		{"TXT", internal.FormatTXT, "txt:Text (encoded):UTF8", "Text (encoded):UTF8"},
		{"DOCX", internal.FormatDOCX, "docx", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, filter := c.getOutputFormatAndFilter(tt.format)
			if format != tt.wantFormat {
				t.Errorf("Format = %s, want %s", format, tt.wantFormat)
			}
			if filter != tt.wantFilterPart {
				t.Errorf("Filter = %s, want %s", filter, tt.wantFilterPart)
			}
		})
	}
}

func TestCheckAvailability(t *testing.T) {
	c := NewConverter("") // Use default PATH lookup
	err := c.CheckAvailability()
	if err != nil {
		t.Skipf("LibreOffice not available: %v", err)
	}
}

func TestGetVersion(t *testing.T) {
	c := NewConverter("")
	version, err := c.GetVersion()
	if err != nil {
		t.Skipf("Could not get LibreOffice version: %v", err)
	}
	if version == "" || version == "unknown" {
		t.Skip("LibreOffice version unknown or not available")
	}
	t.Logf("LibreOffice version: %s", version)
}

// TestIntegration_MarkdownToPDF tests MD → PDF conversion
func TestIntegration_MarkdownToPDF(t *testing.T) {
	c := NewConverter("")
	if err := c.CheckAvailability(); err != nil {
		t.Skipf("LibreOffice not available: %v", err)
	}

	// Create temp directory
	tmpDir := t.TempDir()

	// Create test Markdown file
	inputFile := filepath.Join(tmpDir, "test.md")
	mdContent := `# Test Document

This is a **test** document with:

- List item 1
- List item 2

## Section 2

Some more content.
`
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Convert to PDF
	outputFile := filepath.Join(tmpDir, "test.pdf")
	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatMD,
		OutputFormat: internal.FormatPDF,
	}

	ctx := context.Background()
	err := c.Convert(ctx, inputFile, outputFile, opts)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Verify output file exists and has content
	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Output file not found: %v", err)
	}
	if stat.Size() == 0 {
		t.Error("Output PDF file is empty")
	}

	t.Logf("Output PDF size: %d bytes", stat.Size())
}

// TestIntegration_DOCXToHTML tests DOCX → HTML conversion
func TestIntegration_DOCXToHTML(t *testing.T) {
	c := NewConverter("")
	if err := c.CheckAvailability(); err != nil {
		t.Skipf("LibreOffice not available: %v", err)
	}

	// Create temp directory
	tmpDir := t.TempDir()

	// First create a DOCX from Markdown using LibreOffice
	inputMD := filepath.Join(tmpDir, "test.md")
	mdContent := `# Test Document

This is a test with **bold** and *italic*.
`
	if err := os.WriteFile(inputMD, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tempDOCX := filepath.Join(tmpDir, "test.docx")
	opts1 := internal.ConversionOptions{
		InputFormat:  internal.FormatMD,
		OutputFormat: internal.FormatDOCX,
	}

	ctx := context.Background()
	if err := c.Convert(ctx, inputMD, tempDOCX, opts1); err != nil {
		t.Fatalf("Failed to create DOCX: %v", err)
	}

	// Now convert DOCX → HTML
	outputHTML := filepath.Join(tmpDir, "output.html")
	opts2 := internal.ConversionOptions{
		InputFormat:  internal.FormatDOCX,
		OutputFormat: internal.FormatHTML,
	}

	err := c.Convert(ctx, tempDOCX, outputHTML, opts2)
	if err != nil {
		t.Fatalf("DOCX to HTML conversion failed: %v", err)
	}

	// Verify output
	stat, err := os.Stat(outputHTML)
	if err != nil {
		t.Fatalf("Output HTML not found: %v", err)
	}
	if stat.Size() == 0 {
		t.Error("Output HTML is empty")
	}

	t.Logf("Successfully converted DOCX to HTML (%d bytes)", stat.Size())
}
