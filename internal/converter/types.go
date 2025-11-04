package converter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
	"github.com/valpere/yakateka/internal/converter/config"
	"github.com/valpere/yakateka/internal/converter/generic"
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

// LoadFromConfig loads converters from configuration
func (f *Factory) LoadFromConfig(cfg *config.ConverterConfig) error {
	for name, tool := range cfg.Converters {
		// Skip plaintext - it's handled specially
		if name == "plaintext" {
			continue
		}

		converter := generic.NewConverter(name, tool, cfg.Profiles)
		f.Register(name, converter)

		log.Debug().
			Str("converter", name).
			Strs("input", tool.Formats.Input).
			Strs("output", tool.Formats.Output).
			Msg("Registered converter from config")
	}
	return nil
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

	// Try direct 2-step conversion
	for _, intermediate := range intermediateFormats {
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

	// Try multi-step conversion (up to 4 steps)
	pipeline, err := f.tryMultiStepPipeline(from, to, intermediateFormats, 4)
	if err == nil {
		// Build via description
		var via []string
		for _, step := range pipeline {
			via = append(via, string(step.ToFormat))
		}
		log.Info().
			Str("from", string(from)).
			Str("to", string(to)).
			Str("via", strings.Join(via[:len(via)-1], "→")).
			Int("steps", len(pipeline)).
			Msg("Built multi-step conversion pipeline")
		return pipeline, nil
	}

	return nil, fmt.Errorf("no conversion pipeline found for %s → %s (requires additional converters)", from, to)
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

// tryMultiStepPipeline attempts to build a multi-step pipeline using BFS
// maxSteps limits the search depth (e.g., 4 for up to 4-step pipelines)
func (f *Factory) tryMultiStepPipeline(from, to internal.DocumentFormat, intermediates []internal.DocumentFormat, maxSteps int) ([]ConversionStep, error) {
	// Use breadth-first search to find shortest path
	type node struct {
		format   internal.DocumentFormat
		path     []ConversionStep
		distance int
	}

	queue := []node{{format: from, path: nil, distance: 0}}
	visited := make(map[internal.DocumentFormat]bool)
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Stop if we've exceeded max steps
		if current.distance >= maxSteps {
			continue
		}

		// Try all possible next formats (intermediates + target)
		candidates := append([]internal.DocumentFormat{to}, intermediates...)
		for _, nextFormat := range candidates {
			// Skip if already visited
			if visited[nextFormat] {
				continue
			}

			// Try to find converter for current → next
			converter, err := f.GetConverter(current.format, nextFormat)
			if err != nil {
				continue // No converter available
			}

			// Build new path
			newPath := make([]ConversionStep, len(current.path))
			copy(newPath, current.path)
			newPath = append(newPath, ConversionStep{
				FromFormat: current.format,
				ToFormat:   nextFormat,
				Converter:  converter,
			})

			// Check if we reached the target
			if nextFormat == to {
				return newPath, nil
			}

			// Add to queue for further exploration
			visited[nextFormat] = true
			queue = append(queue, node{
				format:   nextFormat,
				path:     newPath,
				distance: current.distance + 1,
			})
		}
	}

	return nil, fmt.Errorf("no multi-step pipeline found")
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
