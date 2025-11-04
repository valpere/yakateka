package converter

import (
	"context"

	"github.com/valpere/yakateka/internal"
)

// Factory creates converters based on input/output formats
type Factory struct {
	converters map[string]internal.Converter
}

// NewFactory creates a new converter factory
func NewFactory() *Factory {
	return &Factory{
		converters: make(map[string]internal.Converter),
	}
}

// Register registers a converter for specific formats
func (f *Factory) Register(name string, converter internal.Converter) {
	f.converters[name] = converter
}

// GetConverter returns a converter that supports the given formats
func (f *Factory) GetConverter(inputFormat, outputFormat internal.DocumentFormat) (internal.Converter, error) {
	// Try to find a converter that supports both formats
	for _, converter := range f.converters {
		canConvert := false
		for _, inFmt := range converter.SupportedInputFormats() {
			if inFmt == inputFormat {
				for _, outFmt := range converter.SupportedOutputFormats() {
					if outFmt == outputFormat {
						canConvert = true
						break
					}
				}
				break
			}
		}
		if canConvert {
			return converter, nil
		}
	}

	return nil, internal.ErrUnsupportedConversion
}

// Convert performs document conversion using the appropriate converter
func (f *Factory) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	converter, err := f.GetConverter(opts.InputFormat, opts.OutputFormat)
	if err != nil {
		return err
	}

	return converter.Convert(ctx, input, output, opts)
}
