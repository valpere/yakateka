package libreoffice

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles document conversions using LibreOffice
type Converter struct {
	sofficePath string // Path to soffice/libreoffice binary
}

// NewConverter creates a new LibreOffice converter
func NewConverter(sofficePath string) *Converter {
	if sofficePath == "" {
		sofficePath = "soffice" // Use PATH
	}
	return &Converter{
		sofficePath: sofficePath,
	}
}

// SupportedInputFormats returns formats this converter can read
// LibreOffice can read a wide variety of document formats
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatPDF,   // PDF (limited support - mainly for export)
		internal.FormatDOC,   // Microsoft Word 97-2003
		internal.FormatDOCX,  // Microsoft Word 2007+
		internal.FormatODT,   // OpenDocument Text
		internal.FormatRTF,   // Rich Text Format
		internal.FormatPS,    // PostScript
	}
}

// SupportedOutputFormats returns formats this converter can write
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatPDF,   // PDF export
		internal.FormatHTML,  // HTML export (preserves structure)
		internal.FormatTXT,   // Plain text export (loses structure - use HTML instead)
		internal.FormatDOCX,  // Microsoft Word 2007+
		internal.FormatODT,   // OpenDocument Text
		internal.FormatRTF,   // Rich Text Format
		// NOTE: MD removed - LibreOffice doesn't support Markdown natively
		// Use HTML → MD via Pandoc instead to preserve structure
	}
}
