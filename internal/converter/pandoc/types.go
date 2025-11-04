package pandoc

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles document conversions using Pandoc
type Converter struct {
	pandocPath string // Path to pandoc binary
	extraArgs  []string
}

// NewConverter creates a new Pandoc converter
func NewConverter(pandocPath string, extraArgs []string) *Converter {
	if pandocPath == "" {
		pandocPath = "pandoc" // Use PATH
	}
	return &Converter{
		pandocPath: pandocPath,
		extraArgs:  extraArgs,
	}
}

// SupportedInputFormats returns formats this converter can read
// NOTE: Pandoc cannot read PDF or plain TXT files - it can only write to them
// Based on: pandoc --list-input-formats
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatDOCX,
		internal.FormatODT,
		internal.FormatRTF,
		internal.FormatHTML,
		internal.FormatMD,
		internal.FormatEPUB,
		internal.FormatFB2,
		// Note: TXT (plain text) is NOT in this list - Pandoc cannot read it
		internal.FormatJSON,
		internal.FormatCSV,
		internal.FormatLaTeX,
		internal.FormatRST,
	}
}

// SupportedOutputFormats returns formats this converter can write
// Based on: pandoc --list-output-formats
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatTXT,
		internal.FormatMD,
		internal.FormatHTML,
		internal.FormatDOCX,
		internal.FormatODT,
		internal.FormatPDF,
		internal.FormatRTF,
		internal.FormatEPUB,
		internal.FormatFB2,
		internal.FormatJSON,
		internal.FormatCSV,
		internal.FormatLaTeX,
		internal.FormatRST,
	}
}
