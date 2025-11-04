package internal

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	ErrUnsupportedConversion = errors.New("unsupported conversion")
	ErrUnsupportedFormat     = errors.New("unsupported format")
	ErrInvalidInput          = errors.New("invalid input")
	ErrConversionFailed      = errors.New("conversion failed")
)

// DocumentFormat represents a document format type
type DocumentFormat string

const (
	// Document formats
	FormatPDF  DocumentFormat = "pdf"
	FormatEPUB DocumentFormat = "epub"
	FormatFB2  DocumentFormat = "fb2"
	FormatDJVU DocumentFormat = "djvu"
	FormatMOBI DocumentFormat = "mobi"
	FormatDOCX DocumentFormat = "docx"
	FormatDOC  DocumentFormat = "doc"
	FormatODT  DocumentFormat = "odt"
	FormatRTF  DocumentFormat = "rtf"
	FormatTXT  DocumentFormat = "txt"
	FormatHTML DocumentFormat = "html"
	FormatMD   DocumentFormat = "md"
	FormatJSON DocumentFormat = "json"
	FormatYAML DocumentFormat = "yaml"

	// Image formats
	FormatPNG  DocumentFormat = "png"
	FormatJPG  DocumentFormat = "jpg"
	FormatJPEG DocumentFormat = "jpeg"
	FormatTIFF DocumentFormat = "tiff"
	FormatBMP  DocumentFormat = "bmp"
	FormatWEBP DocumentFormat = "webp"
)

// Config represents the global application configuration
type Config struct {
	OCR       OCRConfig       `mapstructure:"ocr"`
	Converter ConverterConfig `mapstructure:"converter"`
	Metadata  MetadataConfig  `mapstructure:"metadata"`
	Output    OutputConfig    `mapstructure:"output"`
	Logging   LoggingConfig   `mapstructure:"logging"`
}

// OCRConfig represents OCR-related configuration
type OCRConfig struct {
	Engine     string           `mapstructure:"engine"`
	Languages  []string         `mapstructure:"languages"`
	DPI        int              `mapstructure:"dpi"`
	Preprocess PreprocessConfig `mapstructure:"preprocess"`
}

// PreprocessConfig represents image preprocessing configuration
type PreprocessConfig struct {
	Denoise   bool   `mapstructure:"denoise"`
	Deskew    bool   `mapstructure:"deskew"`
	Threshold string `mapstructure:"threshold"`
}

// ConverterConfig represents converter-related configuration
type ConverterConfig struct {
	PDF    PDFConfig    `mapstructure:"pdf"`
	Pandoc PandocConfig `mapstructure:"pandoc"`
	Image  ImageConfig  `mapstructure:"image"`
}

// PDFConfig represents PDF converter configuration
type PDFConfig struct {
	Engine  string `mapstructure:"engine"`
	Quality string `mapstructure:"quality"`
}

// PandocConfig represents Pandoc converter configuration
type PandocConfig struct {
	Path      string   `mapstructure:"path"`
	ExtraArgs []string `mapstructure:"extra_args"`
}

// ImageConfig represents image converter configuration
type ImageConfig struct {
	Library string `mapstructure:"library"`
	Format  string `mapstructure:"format"`
	DPI     int    `mapstructure:"dpi"`
}

// MetadataConfig represents metadata handling configuration
type MetadataConfig struct {
	Checksum string `mapstructure:"checksum"`
	Sidecar  bool   `mapstructure:"sidecar"`
	Embed    bool   `mapstructure:"embed"`
}

// OutputConfig represents output formatting configuration
type OutputConfig struct {
	Format string `mapstructure:"format"`
	Pretty bool   `mapstructure:"pretty"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// DocumentMetadata represents metadata extracted from or to be written to a document
type DocumentMetadata struct {
	Title       string            `json:"title,omitempty"`
	Author      string            `json:"author,omitempty"`
	Language    string            `json:"language,omitempty"`
	Category    string            `json:"category,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
	Created     time.Time         `json:"created,omitempty"`
	Modified    time.Time         `json:"modified,omitempty"`
	PageCount   int               `json:"page_count,omitempty"`
	Checksum    string            `json:"checksum,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

// ConversionOptions represents options for document conversion
type ConversionOptions struct {
	InputFormat  DocumentFormat    `json:"input_format"`
	OutputFormat DocumentFormat    `json:"output_format"`
	Quality      string            `json:"quality,omitempty"`
	DPI          int               `json:"dpi,omitempty"`
	OCR          bool              `json:"ocr,omitempty"`
	OCRLanguages []string          `json:"ocr_languages,omitempty"`
	Via          string            `json:"via,omitempty"` // Converter to use (pandoc, libreoffice, etc.)
	Extra        map[string]string `json:"extra,omitempty"`
}

// ExtractionOptions represents options for content extraction
type ExtractionOptions struct {
	OCR          bool     `json:"ocr"`
	OCRLanguages []string `json:"ocr_languages,omitempty"`
	DPI          int      `json:"dpi,omitempty"`
	ExtractType  string   `json:"extract_type"` // text, tables, formulas, code
	Format       string   `json:"format"`       // Output format
}

// Result represents a generic operation result
type Result struct {
	Success bool              `json:"success"`
	Data    interface{}       `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// Converter is the interface for document format converters
type Converter interface {
	// Convert converts a document from input to output format
	Convert(ctx context.Context, input, output string, opts ConversionOptions) error

	// SupportedInputFormats returns formats this converter can read
	SupportedInputFormats() []DocumentFormat

	// SupportedOutputFormats returns formats this converter can write
	SupportedOutputFormats() []DocumentFormat
}

// Parser is the interface for document parsers
type Parser interface {
	// Parse extracts structure and metadata from a document
	Parse(ctx context.Context, input string) (*DocumentMetadata, error)

	// SupportedFormats returns formats this parser can handle
	SupportedFormats() []DocumentFormat
}

// OCREngine is the interface for OCR engines
type OCREngine interface {
	// ExtractText performs OCR on an image or scanned document
	ExtractText(ctx context.Context, input string, opts ExtractionOptions) (string, error)

	// SupportedLanguages returns languages this engine supports
	SupportedLanguages() []string
}

// Extractor is the interface for content extractors
type Extractor interface {
	// Extract extracts specific content from a document
	Extract(ctx context.Context, input string, opts ExtractionOptions) (interface{}, error)

	// SupportedTypes returns extraction types this extractor supports
	SupportedTypes() []string
}
