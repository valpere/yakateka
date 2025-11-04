package converter

import (
	"context"
	"testing"

	"github.com/valpere/yakateka/internal"
)

// mockConverter is a mock converter for testing
type mockConverter struct {
	inputFormats  []internal.DocumentFormat
	outputFormats []internal.DocumentFormat
}

func (m *mockConverter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	return nil
}

func (m *mockConverter) SupportedInputFormats() []internal.DocumentFormat {
	return m.inputFormats
}

func (m *mockConverter) SupportedOutputFormats() []internal.DocumentFormat {
	return m.outputFormats
}

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	if factory == nil {
		t.Fatal("NewFactory returned nil")
	}

	if factory.converters == nil {
		t.Error("Factory converters map should not be nil")
	}
}

func TestFactoryRegister(t *testing.T) {
	factory := NewFactory()

	mock := &mockConverter{
		inputFormats:  []internal.DocumentFormat{internal.FormatPDF},
		outputFormats: []internal.DocumentFormat{internal.FormatTXT},
	}

	factory.Register("test", mock)

	if len(factory.converters) != 1 {
		t.Errorf("Expected 1 converter, got %d", len(factory.converters))
	}

	if factory.converters["test"] != mock {
		t.Error("Registered converter does not match")
	}
}

func TestFactoryGetConverter(t *testing.T) {
	factory := NewFactory()

	mock := &mockConverter{
		inputFormats:  []internal.DocumentFormat{internal.FormatPDF},
		outputFormats: []internal.DocumentFormat{internal.FormatTXT},
	}

	factory.Register("test", mock)

	// Test successful retrieval
	converter, err := factory.GetConverter(internal.FormatPDF, internal.FormatTXT)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if converter != mock {
		t.Error("Retrieved converter does not match")
	}

	// Test unsupported conversion
	converter, err = factory.GetConverter(internal.FormatPDF, internal.FormatDOCX)
	if err == nil {
		t.Error("Expected error for unsupported conversion")
	}
	if converter != nil {
		t.Error("Expected nil converter for unsupported conversion")
	}
}

func TestFactoryConvert(t *testing.T) {
	factory := NewFactory()

	mock := &mockConverter{
		inputFormats:  []internal.DocumentFormat{internal.FormatPDF},
		outputFormats: []internal.DocumentFormat{internal.FormatTXT},
	}

	factory.Register("test", mock)

	opts := internal.ConversionOptions{
		InputFormat:  internal.FormatPDF,
		OutputFormat: internal.FormatTXT,
	}

	// Test successful conversion
	err := factory.Convert(context.Background(), "input.pdf", "output.txt", opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test unsupported conversion
	opts.OutputFormat = internal.FormatDOCX
	err = factory.Convert(context.Background(), "input.pdf", "output.docx", opts)
	if err == nil {
		t.Error("Expected error for unsupported conversion")
	}
}
