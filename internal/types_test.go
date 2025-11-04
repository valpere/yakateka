package internal

import (
	"testing"
)

func TestDocumentFormats(t *testing.T) {
	// Test document format constants
	formats := []DocumentFormat{
		FormatPDF, FormatEPUB, FormatFB2, FormatDJVU, FormatMOBI,
		FormatDOCX, FormatDOC, FormatODT, FormatRTF, FormatTXT,
		FormatHTML, FormatMD, FormatJSON, FormatYAML,
		FormatPNG, FormatJPG, FormatJPEG, FormatTIFF, FormatBMP, FormatWEBP,
	}

	for _, format := range formats {
		if format == "" {
			t.Errorf("Format should not be empty")
		}
		if len(string(format)) < 2 {
			t.Errorf("Format '%s' seems too short", format)
		}
	}
}

func TestDocumentMetadata(t *testing.T) {
	// Test DocumentMetadata structure
	metadata := DocumentMetadata{
		Title:    "Test Document",
		Author:   "Test Author",
		Language: "en",
		Tags:     []string{"test", "example"},
		Custom:   map[string]string{"key": "value"},
	}

	if metadata.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", metadata.Title)
	}

	if len(metadata.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(metadata.Tags))
	}

	if metadata.Custom["key"] != "value" {
		t.Errorf("Expected custom field 'key' to be 'value', got '%s'", metadata.Custom["key"])
	}
}

func TestConversionOptions(t *testing.T) {
	// Test ConversionOptions structure
	opts := ConversionOptions{
		InputFormat:  FormatPDF,
		OutputFormat: FormatTXT,
		Quality:      "high",
		DPI:          300,
		OCR:          true,
		OCRLanguages: []string{"en", "uk"},
	}

	if opts.InputFormat != FormatPDF {
		t.Errorf("Expected input format PDF, got %s", opts.InputFormat)
	}

	if opts.OutputFormat != FormatTXT {
		t.Errorf("Expected output format TXT, got %s", opts.OutputFormat)
	}

	if opts.DPI != 300 {
		t.Errorf("Expected DPI 300, got %d", opts.DPI)
	}

	if !opts.OCR {
		t.Error("Expected OCR to be enabled")
	}
}

func TestResult(t *testing.T) {
	// Test success result
	successResult := Result{
		Success: true,
		Data:    "test data",
	}

	if !successResult.Success {
		t.Error("Expected success result")
	}

	if successResult.Data != "test data" {
		t.Errorf("Expected data 'test data', got '%v'", successResult.Data)
	}

	// Test error result
	errorResult := Result{
		Success: false,
		Error:   "test error",
	}

	if errorResult.Success {
		t.Error("Expected error result")
	}

	if errorResult.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", errorResult.Error)
	}
}
