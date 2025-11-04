package plaintext

import (
	"context"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Convert converts plain text to HTML or Markdown
func (c *Converter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	// Read input file
	content, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	text := string(content)

	var result string
	switch opts.OutputFormat {
	case internal.FormatHTML:
		result = convertToHTML(text)
	case internal.FormatMD:
		result = convertToMarkdown(text)
	default:
		return fmt.Errorf("%w: plaintext converter only supports HTML and MD output", internal.ErrUnsupportedFormat)
	}

	// Write output
	err = os.WriteFile(output, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Info().
		Str("input", input).
		Str("output", output).
		Str("from", "txt").
		Str("to", string(opts.OutputFormat)).
		Int("bytes", len(result)).
		Msg("Successfully converted plain text")

	return nil
}

// convertToHTML wraps plain text in basic HTML structure
func convertToHTML(text string) string {
	// Escape HTML entities
	escaped := html.EscapeString(text)

	// Convert paragraphs (double newlines)
	paragraphs := strings.Split(escaped, "\n\n")
	var htmlParagraphs []string
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			// Preserve single newlines as <br>
			p = strings.ReplaceAll(p, "\n", "<br>\n")
			htmlParagraphs = append(htmlParagraphs, "<p>"+p+"</p>")
		}
	}

	html := `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Converted from Plain Text</title>
</head>
<body>
` + strings.Join(htmlParagraphs, "\n") + `
</body>
</html>`

	return html
}

// convertToMarkdown converts plain text to Markdown (minimal changes)
func convertToMarkdown(text string) string {
	// For plain text to markdown, we mostly preserve the structure
	// Just ensure proper paragraph breaks
	lines := strings.Split(text, "\n")
	var result []string

	for i, line := range lines {
		// If line is empty and next line exists, add extra newline for paragraph break
		if strings.TrimSpace(line) == "" && i < len(lines)-1 && strings.TrimSpace(lines[i+1]) != "" {
			result = append(result, "")
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
