package pandoc

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/valpere/yakateka/internal"
)

// NOTE: Pandoc CANNOT convert FROM PDF - only TO PDF
// PDF → Text conversion requires LibreOffice or OCR+Ollama approach
// This test is commented out to document this limitation

// TestIntegration_EPUBToText tests EPUB to text conversion with real files
func TestIntegration_EPUBToText(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping integration test")
	}

	// Check if test documents directory exists
	// Try multiple possible paths
	testDocsDirs := []string{
		"../../../library4tests",
		"/home/val/wrk/projects/library/library4tests",
	}

	var testDocsDir string
	for _, dir := range testDocsDirs {
		if _, err := os.Stat(dir); err == nil {
			testDocsDir = dir
			break
		}
	}

	if testDocsDir == "" {
		t.Skip("Test documents directory not found, skipping integration test")
	}

	// Create temp directory for output
	tmpDir, err := os.MkdirTemp("", "yakateka-pandoc-epub-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with NoSQL_Distilled-2.epub
	inputEPUB := filepath.Join(testDocsDir, "NoSQL_Distilled-2.epub")
	if _, err := os.Stat(inputEPUB); os.IsNotExist(err) {
		t.Skipf("Test EPUB not found: %s", inputEPUB)
	}

	outputTXT := filepath.Join(tmpDir, "NoSQL_Distilled.txt")

	// Create converter
	c := NewConverter(pandocPath, nil)

	// Convert EPUB to text
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatEPUB,
		OutputFormat: internal.FormatTXT,
	}

	t.Logf("Converting %s to %s", inputEPUB, outputTXT)
	err = c.Convert(ctx, inputEPUB, outputTXT, opts)
	if err != nil {
		t.Fatalf("Convert() failed: %v", err)
	}

	// Verify output file was created
	stat, err := os.Stat(outputTXT)
	if os.IsNotExist(err) {
		t.Fatal("output file was not created")
	}
	if err != nil {
		t.Fatalf("failed to stat output file: %v", err)
	}

	// Verify output has content
	if stat.Size() == 0 {
		t.Fatal("output file is empty")
	}

	t.Logf("Output file size: %d bytes", stat.Size())

	// Read and verify output contains expected text
	content, err := os.ReadFile(outputTXT)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	text := string(content)
	if len(text) == 0 {
		t.Fatal("output text is empty")
	}

	// Verify it contains expected content
	expectedStrings := []string{
		"NoSQL", // Book title
		"database", // Common word in database books
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(strings.ToLower(text), strings.ToLower(expected)) {
			t.Logf("WARNING: Expected text %q not found in output", expected)
			t.Logf("First 500 chars of output:\n%s", text[:min(500, len(text))])
		}
	}
}

// TestIntegration_MarkdownToPDF tests Markdown to PDF conversion
func TestIntegration_MarkdownToPDF(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping integration test")
	}

	// Create temp directory for test files
	tmpDir, err := os.MkdirTemp("", "yakateka-pandoc-md2pdf-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test markdown file
	inputMD := filepath.Join(tmpDir, "test.md")
	markdown := `# Test Document

This is a test document for PDF generation.

## Features

- Bullet points
- **Bold text**
- *Italic text*

## Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

## Conclusion

This tests Pandoc's Markdown to PDF conversion.
`
	if err := os.WriteFile(inputMD, []byte(markdown), 0644); err != nil {
		t.Fatalf("failed to write test markdown: %v", err)
	}

	outputPDF := filepath.Join(tmpDir, "test.pdf")

	// Create converter
	c := NewConverter(pandocPath, nil)

	// Convert Markdown to PDF
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatMD,
		OutputFormat: internal.FormatPDF,
		Quality:      "medium",
	}

	t.Logf("Converting %s to %s", inputMD, outputPDF)
	err = c.Convert(ctx, inputMD, outputPDF, opts)
	if err != nil {
		t.Fatalf("Convert() failed: %v", err)
	}

	// Verify output file was created
	stat, err := os.Stat(outputPDF)
	if os.IsNotExist(err) {
		t.Fatal("output PDF was not created")
	}
	if err != nil {
		t.Fatalf("failed to stat output PDF: %v", err)
	}

	// Verify output has content
	if stat.Size() == 0 {
		t.Fatal("output PDF is empty")
	}

	t.Logf("Output PDF size: %d bytes", stat.Size())

	// Verify it's a valid PDF (starts with %PDF)
	header := make([]byte, 4)
	f, err := os.Open(outputPDF)
	if err != nil {
		t.Fatalf("failed to open output PDF: %v", err)
	}
	defer f.Close()

	_, err = f.Read(header)
	if err != nil {
		t.Fatalf("failed to read PDF header: %v", err)
	}

	if string(header) != "%PDF" {
		t.Errorf("output is not a valid PDF, header: %q", string(header))
	}
}

// TestIntegration_DOCXToText tests DOCX to text conversion
func TestIntegration_DOCXToText(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping integration test")
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "yakateka-pandoc-docx-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// First create a DOCX from markdown (we need a test DOCX)
	inputMD := filepath.Join(tmpDir, "test.md")
	markdown := `# Test DOCX Conversion

This is a test document.

## Section 1

Some text with **bold** and *italic*.

## Section 2

More content here.
`
	if err := os.WriteFile(inputMD, []byte(markdown), 0644); err != nil {
		t.Fatalf("failed to write test markdown: %v", err)
	}

	testDOCX := filepath.Join(tmpDir, "test.docx")
	c := NewConverter(pandocPath, nil)

	// Convert MD to DOCX first
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	err = c.Convert(ctx1, inputMD, testDOCX, internal.ConversionOptions{
		InputFormat:  internal.FormatMD,
		OutputFormat: internal.FormatDOCX,
	})
	if err != nil {
		t.Fatalf("failed to create test DOCX: %v", err)
	}

	// Now convert DOCX to text
	outputTXT := filepath.Join(tmpDir, "output.txt")

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatDOCX,
		OutputFormat: internal.FormatTXT,
	}

	t.Logf("Converting %s to %s", testDOCX, outputTXT)
	err = c.Convert(ctx2, testDOCX, outputTXT, opts)
	if err != nil {
		t.Fatalf("Convert() failed: %v", err)
	}

	// Verify output
	stat, err := os.Stat(outputTXT)
	if os.IsNotExist(err) {
		t.Fatal("output text file was not created")
	}
	if stat.Size() == 0 {
		t.Fatal("output text file is empty")
	}

	// Read and verify content
	content, err := os.ReadFile(outputTXT)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	text := string(content)
	expectedStrings := []string{"Test DOCX Conversion", "Section 1", "Section 2"}
	for _, expected := range expectedStrings {
		if !strings.Contains(text, expected) {
			t.Errorf("Expected text %q not found in output", expected)
		}
	}

	t.Logf("Successfully converted DOCX to text (%d bytes)", stat.Size())
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
