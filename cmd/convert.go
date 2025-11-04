package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valpere/yakateka/internal"
	"github.com/valpere/yakateka/internal/converter"
	"github.com/valpere/yakateka/internal/converter/calibre"
	"github.com/valpere/yakateka/internal/converter/djvu"
	"github.com/valpere/yakateka/internal/converter/libreoffice"
	"github.com/valpere/yakateka/internal/converter/pandoc"
	"github.com/valpere/yakateka/internal/converter/plaintext"
	"github.com/valpere/yakateka/internal/converter/postscript"
)

var (
	inputFormat  string
	outputFormat string
	quality      string
	dpi          int
	via          string
	timeout      int
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <input> <output>",
	Short: "Convert document between formats",
	Long: `Convert documents between different formats.

Supported conversions:
  - PDF → TXT (text extraction)
  - PDF → PNG/JPG (coming soon)
  - More formats coming in future phases

Examples:
  # Convert PDF to text (auto-detect formats from extensions)
  yakateka convert document.pdf document.txt

  # Explicitly specify formats
  yakateka convert document.pdf output.txt --from pdf --to txt

  # Use specific converter
  yakateka convert notes.md document.pdf --via pandoc`,
	Args: cobra.ExactArgs(2),
	RunE: runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Flags
	convertCmd.Flags().StringVarP(&inputFormat, "from", "f", "",
		"input format (auto-detected from extension if not specified)")
	convertCmd.Flags().StringVarP(&outputFormat, "to", "t", "",
		"output format (auto-detected from extension if not specified)")
	convertCmd.Flags().StringVar(&quality, "quality", "",
		"conversion quality (low, medium, high)")
	convertCmd.Flags().IntVar(&dpi, "dpi", 0,
		"DPI for image conversions (default from config)")
	convertCmd.Flags().StringVar(&via, "via", "",
		"specific converter to use (pandoc, libreoffice, etc.)")
	convertCmd.Flags().IntVar(&timeout, "timeout", 300,
		"conversion timeout in seconds (default 300 = 5 minutes)")
}

func runConvert(cmd *cobra.Command, args []string) error {
	input := args[0]
	output := args[1]

	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", input)
	}

	// Auto-detect formats from extensions if not specified
	if inputFormat == "" {
		inputFormat = strings.TrimPrefix(filepath.Ext(input), ".")
		if inputFormat == "" {
			return fmt.Errorf("cannot detect input format, please specify with --from")
		}
	}

	if outputFormat == "" {
		outputFormat = strings.TrimPrefix(filepath.Ext(output), ".")
		if outputFormat == "" {
			return fmt.Errorf("cannot detect output format, please specify with --to")
		}
	}

	// Normalize formats to lowercase
	inputFormat = strings.ToLower(inputFormat)
	outputFormat = strings.ToLower(outputFormat)

	log.Info().
		Str("input", input).
		Str("output", output).
		Str("from", inputFormat).
		Str("to", outputFormat).
		Msg("Starting conversion")

	// Build conversion options
	opts := internal.ConversionOptions{
		InputFormat:  internal.DocumentFormat(inputFormat),
		OutputFormat: internal.DocumentFormat(outputFormat),
		Quality:      quality,
		DPI:          dpi,
		Via:          via,
	}

	// Use quality from config if not specified
	if opts.Quality == "" {
		opts.Quality = viper.GetString("converter.pdf.quality")
	}

	// Use DPI from config if not specified
	if opts.DPI == 0 {
		opts.DPI = viper.GetInt("converter.pdf.dpi")
	}

	// Use timeout from config if not specified (default flag value is 300)
	if cmd.Flags().Changed("timeout") {
		// User explicitly set timeout via flag, use it
	} else {
		// Use timeout from config
		timeout = viper.GetInt("converter.timeout")
	}

	// Create converter factory
	factory := converter.NewFactory()

	// Register Pandoc converter
	pandocPath := viper.GetString("converter.pandoc.path")
	pandocExtraArgs := viper.GetStringSlice("converter.pandoc.extra_args")
	pandocConverter := pandoc.NewConverter(pandocPath, pandocExtraArgs)
	factory.Register("pandoc", pandocConverter)

	// Register DjVu converter
	djvutxtPath := viper.GetString("converter.djvu.djvutxt_path")
	djvupsPath := viper.GetString("converter.djvu.djvups_path")
	djvuConverter := djvu.NewConverter(djvutxtPath, djvupsPath)
	factory.Register("djvu", djvuConverter)

	// Register PlainText converter (for TXT → HTML/MD pipeline support)
	plaintextConverter := plaintext.NewConverter()
	factory.Register("plaintext", plaintextConverter)

	// Register PostScript converter (for PS → PDF pipeline support)
	ps2pdfPath := viper.GetString("converter.postscript.ps2pdf_path")
	psConverter := postscript.NewConverter(ps2pdfPath)
	factory.Register("postscript", psConverter)

	// Register LibreOffice converter (for PDF/DOC → HTML/TXT conversion)
	sofficePath := viper.GetString("converter.libreoffice.soffice_path")
	libreofficeConverter := libreoffice.NewConverter(sofficePath)
	factory.Register("libreoffice", libreofficeConverter)

	// Register Calibre converter (for ebook formats: MOBI, EPUB, FB2, etc.)
	ebookConvertPath := viper.GetString("converter.calibre.ebook_convert_path")
	calibreConverter := calibre.NewConverter(ebookConvertPath)
	factory.Register("calibre", calibreConverter)

	// TODO: Register additional converters:
	// - OCR + Ollama converter

	// Perform conversion with timeout
	timeoutDuration := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	startTime := time.Now()
	err := factory.Convert(ctx, input, output, opts)
	duration := time.Since(startTime)

	if err != nil {
		log.Error().
			Err(err).
			Str("input", input).
			Str("output", output).
			Dur("duration", duration).
			Msg("Conversion failed")
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Get output file size
	stat, _ := os.Stat(output)
	var fileSize int64
	if stat != nil {
		fileSize = stat.Size()
	}

	log.Info().
		Str("output", output).
		Int64("size", fileSize).
		Dur("duration", duration).
		Msg("Conversion completed successfully")

	fmt.Printf("✓ Converted %s → %s (%d bytes) in %v\n",
		input, output, fileSize, duration.Round(time.Millisecond))

	return nil
}
