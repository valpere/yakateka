package pandoc

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/valpere/yakateka/internal"
)

// TestNewConverter verifies converter initialization
func TestNewConverter(t *testing.T) {
	tests := []struct {
		name       string
		pandocPath string
		extraArgs  []string
		wantPath   string
	}{
		{
			name:       "custom path",
			pandocPath: "/usr/local/bin/pandoc",
			extraArgs:  []string{"--verbose"},
			wantPath:   "/usr/local/bin/pandoc",
		},
		{
			name:       "empty path defaults to PATH",
			pandocPath: "",
			extraArgs:  nil,
			wantPath:   "pandoc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter(tt.pandocPath, tt.extraArgs)
			if c.pandocPath != tt.wantPath {
				t.Errorf("pandocPath = %q, want %q", c.pandocPath, tt.wantPath)
			}
			if len(c.extraArgs) != len(tt.extraArgs) {
				t.Errorf("extraArgs length = %d, want %d", len(c.extraArgs), len(tt.extraArgs))
			}
		})
	}
}

// TestSupportedFormats verifies format support
func TestSupportedFormats(t *testing.T) {
	c := NewConverter("", nil)

	inputFormats := c.SupportedInputFormats()
	if len(inputFormats) == 0 {
		t.Error("expected non-empty input formats")
	}

	// Check for key formats
	hasFormat := func(formats []internal.DocumentFormat, format internal.DocumentFormat) bool {
		for _, f := range formats {
			if f == format {
				return true
			}
		}
		return false
	}

	if !hasFormat(inputFormats, internal.FormatDOCX) {
		t.Error("expected DOCX in input formats")
	}
	if !hasFormat(inputFormats, internal.FormatMD) {
		t.Error("expected MD in input formats")
	}
	if !hasFormat(inputFormats, internal.FormatEPUB) {
		t.Error("expected EPUB in input formats")
	}

	// Pandoc cannot read PDF - only write
	if hasFormat(inputFormats, internal.FormatPDF) {
		t.Error("PDF should NOT be in input formats - Pandoc can only write PDF, not read it")
	}

	outputFormats := c.SupportedOutputFormats()
	if len(outputFormats) == 0 {
		t.Error("expected non-empty output formats")
	}

	if !hasFormat(outputFormats, internal.FormatTXT) {
		t.Error("expected TXT in output formats")
	}
	if !hasFormat(outputFormats, internal.FormatHTML) {
		t.Error("expected HTML in output formats")
	}
}

// TestBuildArgs verifies CLI argument construction
func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		output   string
		opts     internal.ConversionOptions
		wantArgs []string
	}{
		{
			name:   "minimal args",
			input:  "input.md",
			output: "output.html",
			opts:   internal.ConversionOptions{},
			wantArgs: []string{
				"input.md",
				"-o", "output.html",
			},
		},
		{
			name:   "with formats",
			input:  "input.pdf",
			output: "output.txt",
			opts: internal.ConversionOptions{
				InputFormat:  internal.FormatPDF,
				OutputFormat: internal.FormatTXT,
			},
			wantArgs: []string{
				"input.pdf",
				"-o", "output.txt",
				"--from", "pdf",
				"--to", "plain",
			},
		},
		{
			name:   "PDF output with high quality",
			input:  "input.md",
			output: "output.pdf",
			opts: internal.ConversionOptions{
				OutputFormat: internal.FormatPDF,
				Quality:      "high",
			},
			wantArgs: []string{
				"input.md",
				"-o", "output.pdf",
				"--to", "pdf",
				"--pdf-engine=xelatex",
			},
		},
		{
			name:   "with extra args from config",
			input:  "input.md",
			output: "output.pdf",
			opts: internal.ConversionOptions{
				Extra: map[string]string{
					"template": "custom.tex",
					"verbose":  "",
				},
			},
			wantArgs: []string{
				"input.md",
				"-o", "output.pdf",
				"--template=custom.tex",
				"--verbose=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter("", nil)
			args := c.buildArgs(tt.input, tt.output, tt.opts)

			// Check that all expected args are present
			for _, want := range tt.wantArgs {
				found := false
				for _, got := range args {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected arg %q not found in %v", want, args)
				}
			}
		})
	}
}

// TestCheckAvailability verifies pandoc availability check
func TestCheckAvailability(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping availability test")
	}

	c := NewConverter(pandocPath, nil)
	if err := c.CheckAvailability(); err != nil {
		t.Errorf("CheckAvailability() failed: %v", err)
	}

	// Test with invalid path
	c = NewConverter("/nonexistent/pandoc", nil)
	if err := c.CheckAvailability(); err == nil {
		t.Error("expected error for nonexistent pandoc path")
	}
}

// TestGetVersion verifies version retrieval
func TestGetVersion(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping version test")
	}

	c := NewConverter(pandocPath, nil)
	version, err := c.GetVersion()
	if err != nil {
		t.Errorf("GetVersion() failed: %v", err)
	}
	if version == "" {
		t.Error("expected non-empty version string")
	}
	t.Logf("Pandoc version: %s", version)
}

// TestConvertMarkdownToHTML tests basic markdown conversion
func TestConvertMarkdownToHTML(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping conversion test")
	}

	// Create temp directory for test files
	tmpDir, err := os.MkdirTemp("", "yakateka-pandoc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test markdown file
	inputPath := filepath.Join(tmpDir, "test.md")
	content := []byte("# Test Document\n\nThis is a **test** with *formatting*.\n")
	if err := os.WriteFile(inputPath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Convert to HTML
	outputPath := filepath.Join(tmpDir, "test.html")
	c := NewConverter(pandocPath, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatMD,
		OutputFormat: internal.FormatHTML,
	}

	if err := c.Convert(ctx, inputPath, outputPath, opts); err != nil {
		t.Fatalf("Convert() failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("output file was not created")
	}

	// Read and verify output contains expected HTML
	htmlContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	html := string(htmlContent)
	if len(html) == 0 {
		t.Error("output HTML is empty")
	}

	// Basic sanity checks
	expectedStrings := []string{"<h1", "Test Document", "<strong>", "test", "<em>", "formatting"}
	for _, expected := range expectedStrings {
		if !contains(html, expected) {
			t.Errorf("output HTML does not contain %q", expected)
		}
	}
}

// TestConvertTimeout verifies timeout handling
func TestConvertTimeout(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping timeout test")
	}

	// Create temp directory for test files
	tmpDir, err := os.MkdirTemp("", "yakateka-pandoc-timeout-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	inputPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputPath, []byte("# Test\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "test.html")
	c := NewConverter(pandocPath, nil)

	// Use extremely short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatMD,
		OutputFormat: internal.FormatHTML,
	}

	err = c.Convert(ctx, inputPath, outputPath, opts)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

// TestConvertInvalidInput verifies error handling for missing input
func TestConvertInvalidInput(t *testing.T) {
	// Find pandoc in PATH
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH, skipping invalid input test")
	}

	c := NewConverter(pandocPath, nil)
	ctx := context.Background()

	err = c.Convert(ctx, "/nonexistent/file.md", "/tmp/output.html", internal.ConversionOptions{})
	if err == nil {
		t.Error("expected error for nonexistent input file")
	}

	// Should return ErrInvalidInput
	if err != nil && !contains(err.Error(), "invalid input") {
		t.Logf("got error: %v", err)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr ||
			(len(s) > len(substr) && contains(s[1:], substr)))))
}
