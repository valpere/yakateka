package postscript

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles PostScript conversions using Ghostscript tools
type Converter struct {
	ps2pdfPath string // Path to ps2pdf binary
}

// NewConverter creates a new PostScript converter
func NewConverter(ps2pdfPath string) *Converter {
	if ps2pdfPath == "" {
		ps2pdfPath = "ps2pdf" // Use PATH
	}
	return &Converter{
		ps2pdfPath: ps2pdfPath,
	}
}

// SupportedInputFormats returns formats this converter can read
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatPS,
	}
}

// SupportedOutputFormats returns formats this converter can write
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatPDF,
	}
}
