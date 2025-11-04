package djvu

import (
	"github.com/valpere/yakateka/internal"
)

// Converter handles DjVu document conversions using DjVuLibre tools
type Converter struct {
	djvutxtPath string // Path to djvutxt binary
	djvupsPath  string // Path to djvups binary
}

// NewConverter creates a new DjVu converter
func NewConverter(djvutxtPath, djvupsPath string) *Converter {
	if djvutxtPath == "" {
		djvutxtPath = "djvutxt" // Use PATH
	}
	if djvupsPath == "" {
		djvupsPath = "djvups" // Use PATH
	}
	return &Converter{
		djvutxtPath: djvutxtPath,
		djvupsPath:  djvupsPath,
	}
}

// SupportedInputFormats returns formats this converter can read
// DjVuLibre can only read DJVU files
func (c *Converter) SupportedInputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatDJVU,
	}
}

// SupportedOutputFormats returns formats this converter can write
// DjVuLibre can extract plain text (djvutxt) or PostScript (djvups)
func (c *Converter) SupportedOutputFormats() []internal.DocumentFormat {
	return []internal.DocumentFormat{
		internal.FormatTXT, // via djvutxt
		internal.FormatPS,  // via djvups
	}
}
