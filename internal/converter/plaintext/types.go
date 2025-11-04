package plaintext

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles plain text conversions (simple wrapping in HTML)
type Converter struct{}

// NewConverter creates a new plain text converter
func NewConverter() *Converter {
	return &Converter{}
}

// SupportedInputFormats returns formats this converter can read
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatTXT,
	}
}

// SupportedOutputFormats returns formats this converter can write
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatHTML,
		internal.FormatMD,
	}
}
