package pdf

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles PDF document conversions
type Converter struct {
	engine string // Engine to use (pdfcpu, unipdf)
}

// NewConverter creates a new PDF converter
func NewConverter(engine string) *Converter {
	if engine == "" {
		engine = "pdfcpu"
	}
	return &Converter{
		engine: engine,
	}
}

// SupportedInputFormats returns formats this converter can read
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatPDF,
	}
}

// SupportedOutputFormats returns formats this converter can write
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatTXT,
		internal.FormatPNG,
		internal.FormatJPG,
	}
}
