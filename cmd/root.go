package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	logLevel  string
	logFormat string
	verbose   bool
	version   = "0.1.0" // Version is set via ldflags during build
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "yakateka",
	Version: version,
	Short:   "YakaTeka - Document processing agent",
	Long: `YakaTeka (ЯкаТека) is a document processing agent that handles
document conversion, parsing, OCR, and annotation.

It operates as a standalone CLI tool focused on individual document operations,
supporting formats like PDF, EPUB, FB2, DJVU, MOBI, DOC/DOCX and more.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return setupLogging()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Command execution failed")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yakateka/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "json", "log format (json, text)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (same as --log-level=debug)")

	// Bind flags to viper
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("log-format"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get home directory")
			os.Exit(1)
		}

		// Search config in home directory with name ".yakateka" (without extension)
		configDir := filepath.Join(home, ".yakateka")
		viper.AddConfigPath(configDir)
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variables
	viper.SetEnvPrefix("YAKATEKA")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Str("config", viper.ConfigFileUsed()).Msg("Using config file")
	}

	// Also try to load converters.yaml from the same directory
	convertersFile := filepath.Join(filepath.Dir(viper.ConfigFileUsed()), "converters.yaml")
	if _, err := os.Stat(convertersFile); err == nil {
		viper.SetConfigFile(convertersFile)
		if err := viper.MergeInConfig(); err != nil {
			log.Warn().Err(err).Str("file", convertersFile).Msg("Failed to merge converters config")
		} else {
			log.Debug().Str("config", convertersFile).Msg("Merged converters config")
		}
	}
}

// setDefaults sets default configuration values
func setDefaults() {
	// OCR defaults
	viper.SetDefault("ocr.engine", "tesseract")
	viper.SetDefault("ocr.languages", []string{"uk", "en", "ru"})
	viper.SetDefault("ocr.dpi", 300)
	viper.SetDefault("ocr.preprocess.denoise", true)
	viper.SetDefault("ocr.preprocess.deskew", true)
	viper.SetDefault("ocr.preprocess.threshold", "auto")

	// Converter defaults
	viper.SetDefault("converter.timeout", 300)
	viper.SetDefault("converter.pdf.engine", "pdfcpu")
	viper.SetDefault("converter.pdf.quality", "high")
	viper.SetDefault("converter.pandoc.path", "/usr/bin/pandoc")
	viper.SetDefault("converter.pandoc.extra_args", []string{})
	viper.SetDefault("converter.image.library", "bimg")
	viper.SetDefault("converter.image.format", "png")
	viper.SetDefault("converter.image.dpi", 300)

	// Metadata defaults
	viper.SetDefault("metadata.checksum", "sha256")
	viper.SetDefault("metadata.sidecar", true)
	viper.SetDefault("metadata.embed", true)

	// Output defaults
	viper.SetDefault("output.format", "json")
	viper.SetDefault("output.pretty", true)

	// Logging defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
}

// setupLogging configures zerolog based on flags and config
func setupLogging() error {
	// Set log level
	level := viper.GetString("log.level")
	if verbose {
		level = "debug"
	}

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	// Set log format
	format := viper.GetString("log.format")
	if format == "text" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Debug().
		Str("level", level).
		Str("format", format).
		Msg("Logging initialized")

	return nil
}
