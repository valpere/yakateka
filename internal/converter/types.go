package converter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
)

// Factory creates converters based on input/output formats
type Factory struct {
	converters map[string]internal.Converter
}

// ConversionStep represents one step in a conversion pipeline
type ConversionStep struct {
	FromFormat internal.DocumentFormat
	ToFormat   internal.DocumentFormat
	Converter  internal.Converter
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
// If no direct converter is available, it will attempt a pipeline conversion via HTML
func (f *Factory) Convert(ctx context.Context, input, output string, opts internal.ConversionOptions) error {
	// Try direct conversion first
	converter, err := f.GetConverter(opts.InputFormat, opts.OutputFormat)
	if err == nil {
		log.Debug().
			Str("from", string(opts.InputFormat)).
			Str("to", string(opts.OutputFormat)).
			Msg("Using direct conversion")
		return converter.Convert(ctx, input, output, opts)
	}

	// No direct converter found, try pipeline via HTML
	log.Debug().
		Str("from", string(opts.InputFormat)).
		Str("to", string(opts.OutputFormat)).
		Msg("No direct converter found, attempting pipeline via HTML")

	pipeline, err := f.buildPipeline(opts.InputFormat, opts.OutputFormat)
	if err != nil {
		return err
	}

	return f.executePipeline(ctx, input, output, opts, pipeline)
}

// buildPipeline attempts to build a conversion pipeline via intermediate formats
func (f *Factory) buildPipeline(from, to internal.DocumentFormat) ([]ConversionStep, error) {
	// Try different intermediate formats in order of preference
	// Priority: PDF/PS (preserves structure) > HTML
	// NOTE: TXT is NOT included as it loses all document structure
	intermediateFormats := []internal.DocumentFormat{
		internal.FormatPDF,  // Best: Preserves full document structure
		internal.FormatPS,   // PostScript: Good structure preservation
		internal.FormatHTML, // Good: HTML preserves structure
	}

	for _, intermediate := range intermediateFormats {
		// Try 2-step pipeline: from → intermediate → to
		pipeline, err := f.try2StepPipeline(from, to, intermediate)
		if err == nil {
			log.Info().
				Str("from", string(from)).
				Str("to", string(to)).
				Str("via", string(intermediate)).
				Int("steps", len(pipeline)).
				Msg("Built 2-step conversion pipeline")
			return pipeline, nil
		}
	}

	return nil, fmt.Errorf("no conversion pipeline found for %s → %s (requires LibreOffice or additional converters)", from, to)
}

// try2StepPipeline attempts to build a 2-step pipeline via an intermediate format
func (f *Factory) try2StepPipeline(from, to, intermediate internal.DocumentFormat) ([]ConversionStep, error) {
	// Step 1: from → intermediate
	converter1, err := f.GetConverter(from, intermediate)
	if err != nil {
		return nil, err
	}

	// Step 2: intermediate → to
	converter2, err := f.GetConverter(intermediate, to)
	if err != nil {
		return nil, err
	}

	return []ConversionStep{
		{FromFormat: from, ToFormat: intermediate, Converter: converter1},
		{FromFormat: intermediate, ToFormat: to, Converter: converter2},
	}, nil
}


// executePipeline executes a multi-step conversion pipeline
func (f *Factory) executePipeline(ctx context.Context, input, output string, opts internal.ConversionOptions, pipeline []ConversionStep) error {
	currentInput := input
	var tempFiles []string

	// Execute each step in the pipeline
	for i, step := range pipeline {
		var currentOutput string

		if i == len(pipeline)-1 {
			// Last step: use final output path
			currentOutput = output
		} else {
			// Intermediate step: create temp file
			tempFile, err := os.CreateTemp("", fmt.Sprintf("yakateka-pipeline-*.%s", step.ToFormat))
			if err != nil {
				return fmt.Errorf("failed to create temp file for pipeline step %d: %w", i+1, err)
			}
			currentOutput = tempFile.Name()
			tempFile.Close()
			tempFiles = append(tempFiles, currentOutput)
		}

		log.Debug().
			Int("step", i+1).
			Int("total", len(pipeline)).
			Str("from", string(step.FromFormat)).
			Str("to", string(step.ToFormat)).
			Str("input", filepath.Base(currentInput)).
			Str("output", filepath.Base(currentOutput)).
			Msg("Executing pipeline step")

		// Create options for this step
		stepOpts := opts
		stepOpts.InputFormat = step.FromFormat
		stepOpts.OutputFormat = step.ToFormat

		// Execute conversion
		err := step.Converter.Convert(ctx, currentInput, currentOutput, stepOpts)
		if err != nil {
			// Clean up temp files on error
			for _, tempFile := range tempFiles {
				os.Remove(tempFile)
			}
			return fmt.Errorf("pipeline step %d failed (%s → %s): %w",
				i+1, step.FromFormat, step.ToFormat, err)
		}

		// Next step's input is current output
		currentInput = currentOutput
	}

	// Clean up intermediate temp files (but not the final output)
	for _, tempFile := range tempFiles {
		log.Debug().Str("file", tempFile).Msg("Removing intermediate temp file")
		os.Remove(tempFile)
	}

	log.Info().
		Int("steps", len(pipeline)).
		Str("input", input).
		Str("output", output).
		Msg("Pipeline conversion completed successfully")

	return nil
}
