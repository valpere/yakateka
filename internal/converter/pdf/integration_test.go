package pdf

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/valpere/yakateka/internal"
)

const (
	testDocumentsDir = "../../../../library4tests"
	testPDF          = "NoSQL_Distilled.pdf"
)

func TestPDFToTextIntegration(t *testing.T) {
	// Skip if test documents directory doesn't exist
	if _, err := os.Stat(testDocumentsDir); os.IsNotExist(err) {
		t.Skip("Test documents directory not found, skipping integration test")
	}

	inputFile := filepath.Join(testDocumentsDir, testPDF)
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Skipf("Test PDF %s not found, skipping integration test", testPDF)
	}

	// Create temporary output directory
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	// Create converter
	converter := NewConverter("pdfcpu")

	// Conversion options
	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatPDF,
		OutputFormat: internal.FormatTXT,
		Quality:      "high",
	}

	// Perform conversion
	ctx := context.Background()
	err := converter.Convert(ctx, inputFile, outputFile, opts)
	if err != nil {
		t.Fatalf("Failed to convert PDF to text: %v", err)
	}

	// Verify output file exists
	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}

	// Verify output file has content
	if stat.Size() == 0 {
		t.Error("Output file is empty")
	}

	// Read output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify we have some text output
	if len(content) > 0 {
		t.Logf("Successfully extracted %d bytes of text from PDF", len(content))
	}

	if len(content) < 100 {
		t.Errorf("Expected more content, got only %d bytes", len(content))
	}

	t.Logf("Integration test successful: converted %s (%d bytes) to text (%d bytes)",
		testPDF, getFileSize(inputFile), len(content))
}

func TestPDFToTextIntegrationWithCLI(t *testing.T) {
	// Skip if test documents directory doesn't exist
	if _, err := os.Stat(testDocumentsDir); os.IsNotExist(err) {
		t.Skip("Test documents directory not found, skipping integration test")
	}

	inputFile := filepath.Join(testDocumentsDir, testPDF)
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Skipf("Test PDF %s not found, skipping integration test", testPDF)
	}

	// This test documents the expected CLI usage
	// Actual CLI testing would use exec.Command in a separate test file

	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "cli_output.txt")

	converter := NewConverter("")
	ctx := context.Background()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatPDF,
		OutputFormat: internal.FormatTXT,
	}

	err := converter.Convert(ctx, inputFile, outputFile, opts)
	if err != nil {
		t.Fatalf("CLI-style conversion failed: %v", err)
	}

	// Verify output
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	t.Logf("CLI integration test successful")
}

// Helper function to get file size
func getFileSize(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.Size()
}
