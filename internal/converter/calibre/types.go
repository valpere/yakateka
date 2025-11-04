package calibre

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles ebook conversions using Calibre's ebook-convert
type Converter struct {
	ebookConvertPath string // Path to ebook-convert binary
}

// NewConverter creates a new Calibre converter
func NewConverter(ebookConvertPath string) *Converter {
	if ebookConvertPath == "" {
		ebookConvertPath = "ebook-convert" // Use PATH
	}
	return &Converter{
		ebookConvertPath: ebookConvertPath,
	}
}

// SupportedInputFormats returns formats this converter can read
// Calibre supports a wide variety of ebook formats
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatMOBI,  // Amazon Kindle
		internal.FormatEPUB,  // EPUB ebooks
		internal.FormatFB2,   // FictionBook 2.0
		internal.FormatHTML,  // HTML files
		internal.FormatTXT,   // Plain text
		internal.FormatPDF,   // PDF (limited support)
		internal.FormatDOCX,  // Microsoft Word
		internal.FormatODT,   // OpenDocument Text
		internal.FormatRTF,   // Rich Text Format
	}
}

// SupportedOutputFormats returns formats this converter can write
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatMOBI,  // Amazon Kindle
		internal.FormatEPUB,  // EPUB ebooks
		internal.FormatFB2,   // FictionBook 2.0
		internal.FormatHTML,  // HTML files
		internal.FormatTXT,   // Plain text
		internal.FormatPDF,   // PDF (via plugins)
		internal.FormatDOCX,  // Microsoft Word
		internal.FormatODT,   // OpenDocument Text
		internal.FormatRTF,   // Rich Text Format
	}
}
