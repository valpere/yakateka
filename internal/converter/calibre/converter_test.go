package calibre

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/valpere/yakateka/internal"
)

func TestNewConverter(t *testing.T) {
	t.Run("custom_path", func(t *testing.T) {
		c := NewConverter("/custom/path/ebook-convert")
		if c.ebookConvertPath != "/custom/path/ebook-convert" {
			t.Errorf("Expected /custom/path/ebook-convert, got %s", c.ebookConvertPath)
		}
	})

	t.Run("empty_path_defaults_to_PATH", func(t *testing.T) {
		c := NewConverter("")
		if c.ebookConvertPath != "ebook-convert" {
			t.Errorf("Expected ebook-convert, got %s", c.ebookConvertPath)
		}
	})
}

func TestSupportedFormats(t *testing.T) {
	c := NewConverter("")

	inputFormats := c.SupportedInputFormats()
	if len(inputFormats) == 0 {
		t.Error("Expected non-zero input formats")
	}

	// Check for key ebook formats
	hasMOBI := false
	hasEPUB := false
	hasFB2 := false
	for _, f := range inputFormats {
		if f == internal.FormatMOBI {
			hasMOBI = true
		}
		if f == internal.FormatEPUB {
			hasEPUB = true
		}
		if f == internal.FormatFB2 {
			hasFB2 = true
		}
	}
	if !hasMOBI {
		t.Error("Expected MOBI in input formats")
	}
	if !hasEPUB {
		t.Error("Expected EPUB in input formats")
	}
	if !hasFB2 {
		t.Error("Expected FB2 in input formats")
	}

	outputFormats := c.SupportedOutputFormats()
	if len(outputFormats) == 0 {
		t.Error("Expected non-zero output formats")
	}
}

func TestCheckAvailability(t *testing.T) {
	c := NewConverter("") // Use default PATH lookup
	err := c.CheckAvailability()
	if err != nil {
		t.Skipf("Calibre ebook-convert not available: %v", err)
	}
}

func TestGetVersion(t *testing.T) {
	c := NewConverter("")
	version, err := c.GetVersion()
	if err != nil {
		t.Skipf("Could not get Calibre version: %v", err)
	}
	if version == "" || version == "unknown" {
		t.Skip("Calibre version unknown or not available")
	}
	t.Logf("Calibre version: %s", version)
}

// TestIntegration_EPUBToMOBI tests EPUB → MOBI conversion
func TestIntegration_EPUBToMOBI(t *testing.T) {
	c := NewConverter("")
	if err := c.CheckAvailability(); err != nil {
		t.Skipf("Calibre not available: %v", err)
	}

	// Use real EPUB file from test library
	inputFile := "/home/val/wrk/projects/library/library4tests/NoSQL_Distilled-2.epub"
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Skipf("Test EPUB file not found: %s", inputFile)
	}

	// Create temp directory for output
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.mobi")

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatEPUB,
		OutputFormat: internal.FormatMOBI,
	}

	ctx := context.Background()
	err := c.Convert(ctx, inputFile, outputFile, opts)
	if err != nil {
		t.Fatalf("EPUB to MOBI conversion failed: %v", err)
	}

	// Verify output file exists and has content
	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Output MOBI file not found: %v", err)
	}
	if stat.Size() == 0 {
		t.Error("Output MOBI file is empty")
	}

	t.Logf("Output MOBI size: %d bytes", stat.Size())
}

// TestIntegration_FB2ToEPUB tests FB2 → EPUB conversion
func TestIntegration_FB2ToEPUB(t *testing.T) {
	c := NewConverter("")
	if err := c.CheckAvailability(); err != nil {
		t.Skipf("Calibre not available: %v", err)
	}

	// Use real FB2 file from test library
	inputFile := "/home/val/wrk/projects/library/library4tests/Breht_Chto-tot-soldat-chto-etot.xxxKvQ.463721.fb2"
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Skipf("Test FB2 file not found: %s", inputFile)
	}

	// Create temp directory for output
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.epub")

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatFB2,
		OutputFormat: internal.FormatEPUB,
		Quality:      "high",
	}

	ctx := context.Background()
	err := c.Convert(ctx, inputFile, outputFile, opts)
	if err != nil {
		t.Fatalf("FB2 to EPUB conversion failed: %v", err)
	}

	// Verify output file exists and has content
	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Output EPUB file not found: %v", err)
	}
	if stat.Size() == 0 {
		t.Error("Output EPUB file is empty")
	}

	t.Logf("Output EPUB size: %d bytes", stat.Size())
}
