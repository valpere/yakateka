package helper

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// HelperConverter uses external helper scripts for conversion
type HelperConverter struct {
	cache    *HelperCache
	executor *Executor
}

// NewHelperConverter creates a converter that uses helper scripts
func NewHelperConverter(cache *HelperCache, executor *Executor) *HelperConverter {
	return &HelperConverter{
		cache:    cache,
		executor: executor,
	}
}

// SupportedInputFormats returns all input formats supported by any helper
func (c *HelperConverter) SupportedInputFormats() []internal.DocumentFormat {
	formatMap := make(map[internal.DocumentFormat]bool)

	for fromFormat := range c.cache.Conversions {
		formatMap[internal.DocumentFormat(fromFormat)] = true
	}

	formats := make([]internal.DocumentFormat, 0, len(formatMap))
	for format := range formatMap {
		formats = append(formats, format)
	}

	return formats
}

// SupportedOutputFormats returns all output formats supported by any helper
func (c *HelperConverter) SupportedOutputFormats() []internal.DocumentFormat {
	formatMap := make(map[internal.DocumentFormat]bool)

	for _, toFormats := range c.cache.Conversions {
		for toFormat := range toFormats {
			formatMap[internal.DocumentFormat(toFormat)] = true
		}
	}

	formats := make([]internal.DocumentFormat, 0, len(formatMap))
	for format := range formatMap {
		formats = append(formats, format)
	}

	return formats
}

// Convert performs document conversion using helper scripts
func (c *HelperConverter) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Validate input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error().Str("input", input).Msg("Input file does not exist")
		return fmt.Errorf("%w: %s", internal.ErrInvalidInput, input)
	}

	// Determine conversion mode from quality setting
	mode := ModeNormal
	switch opts.Quality {
	case "fast":
		mode = ModeFast
	case "high", "quality":
		mode = ModeQuality
	}

	// Find helpers for this conversion
	helpers := c.cache.FindHelpers(opts.InputFormat, opts.OutputFormat, mode)
	if len(helpers) == 0 {
		log.Debug().
			Str("from", string(opts.InputFormat)).
			Str("to", string(opts.OutputFormat)).
			Str("mode", string(mode)).
			Msg("No helpers available for conversion")
		return fmt.Errorf("%w: no helpers support %s → %s",
			internal.ErrUnsupportedConversion, opts.InputFormat, opts.OutputFormat)
	}

	log.Debug().
		Str("from", string(opts.InputFormat)).
		Str("to", string(opts.OutputFormat)).
		Str("mode", string(mode)).
		Int("helpers", len(helpers)).
		Msg("Found helpers for conversion")

	// Try each helper in order
	var lastErr error
	for i, helperEntry := range helpers {
		helperPath := helperEntry.Helper

		log.Info().
			Str("helper", helperPath).
			Int("attempt", i+1).
			Int("total", len(helpers)).
			Str("mode", string(mode)).
			Msg("Attempting conversion with helper")

		err := c.executor.Convert(
			ctx,
			helperPath,
			mode,
			string(opts.InputFormat),
			input,
			string(opts.OutputFormat),
			output,
		)

		if err == nil {
			// Success!
			log.Info().
				Str("helper", helperPath).
				Str("from", string(opts.InputFormat)).
				Str("to", string(opts.OutputFormat)).
				Str("mode", string(mode)).
				Msg("Conversion successful")
			return nil
		}

		// Conversion failed
		log.Warn().
			Err(err).
			Str("helper", helperPath).
			Msg("Helper conversion failed, trying next helper")

		// Mark this helper as failed for this conversion pair
		c.cache.MarkHelperFailed(opts.InputFormat, opts.OutputFormat, helperPath)
		lastErr = err
	}

	// All helpers failed
	log.Error().
		Err(lastErr).
		Str("from", string(opts.InputFormat)).
		Str("to", string(opts.OutputFormat)).
		Int("tried", len(helpers)).
		Msg("All helpers failed")

	return fmt.Errorf("%w: all %d helpers failed, last error: %v",
		internal.ErrConversionFailed, len(helpers), lastErr)
}
