package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valpere/yakateka/internal/helper"
)

var (
	showFormatsMatrix bool
)

// helpersCmd represents the helpers command
var helpersCmd = &cobra.Command{
	Use:   "helpers",
	Short: "Generate helpers cache file",
	Long: `Query all configured helpers and generate helpers.yaml cache file.

This command:
1. Reads helper paths and weights from config
2. Pings each helper to check availability
3. Queries each helper for capabilities (helper.sh info)
4. Generates helpers.yaml with sorted helper lists

The cache file is used at runtime for fast helper lookup.

Use --formats to display a matrix of supported format conversions.`,
	RunE: runHelpers,
}

func init() {
	rootCmd.AddCommand(helpersCmd)
	helpersCmd.Flags().BoolVar(&showFormatsMatrix, "formats", false, "Display format conversion matrix")
}

func runHelpers(cmd *cobra.Command, args []string) error {
	// Load helper configuration
	helperWeights := viper.GetStringMap("helpers.weights")
	if len(helperWeights) == 0 {
		log.Warn().Msg("No helpers configured in config file")
		return fmt.Errorf("no helpers configured (check helpers.weights in config)")
	}

	// Get cache file path
	cacheFile := viper.GetString("helpers.cache_file")
	if cacheFile == "" {
		cacheFile = "helpers.yaml"
	}

	log.Info().
		Int("count", len(helperWeights)).
		Str("cache", cacheFile).
		Msg("Generating helper cache")

	// Create registry
	registry := helper.NewRegistry()

	// Register all helpers
	for path, weight := range helperWeights {
		// Expand environment variables in path (e.g., ${HOME}/path)
		path = os.ExpandEnv(path)

		// Convert relative paths to absolute based on current working directory
		if !filepath.IsAbs(path) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Error().Err(err).Str("helper", path).Msg("Failed to get current working directory for relative helper path")
				return fmt.Errorf("failed to get current working directory for helper %s: %w", path, err)
			}
			path = filepath.Join(cwd, path)
		}
		weightFloat, ok := weight.(float64)
		if !ok {
			log.Warn().
				Str("helper", path).
				Interface("weight", weight).
				Msg("Invalid weight value, skipping helper")
			continue
		}

		registry.Register(path, weightFloat)
		log.Debug().
			Str("helper", path).
			Float64("weight", weightFloat).
			Msg("Registered helper")
	}

	// Initialize helpers (ping + get info)
	executor := helper.NewExecutor(helper.GetTimeout())
	ctx := context.Background()

	if err := registry.Initialize(ctx, executor); err != nil {
		log.Error().Err(err).Msg("Failed to initialize helpers")
		return fmt.Errorf("helper initialization failed: %w", err)
	}

	// Generate cache
	cache, err := registry.GenerateCache()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate cache")
		return fmt.Errorf("cache generation failed: %w", err)
	}

	// Save cache to file
	if err := cache.SaveCache(cacheFile); err != nil {
		log.Error().Err(err).Msg("Failed to save cache")
		return fmt.Errorf("failed to save cache: %w", err)
	}

	// Print summary
	conversionsCount := 0
	for _, toFormats := range cache.Conversions {
		for _, modes := range toFormats {
			conversionsCount += len(modes)
		}
	}

	log.Info().
		Str("cache", cacheFile).
		Int("conversions", conversionsCount).
		Msg("Successfully generated helper cache")

	fmt.Printf("✓ Generated helper cache: %s\n", cacheFile)
	fmt.Printf("  %d conversion paths available\n", conversionsCount)

	// Display format matrix if requested
	if showFormatsMatrix {
		fmt.Println()
		displayFormatMatrix(cache)
	}

	return nil
}

const (
	// minColWidth is the minimum column width for format matrix display
	// Set to 8 to accommodate the "FROM\TO" header (7 chars) plus spacing
	minColWidth = 8
)

// displayFormatMatrix prints a matrix showing supported format conversions
func displayFormatMatrix(cache *helper.HelperCache) {
	// Collect all unique formats that have at least one conversion
	formatSet := make(map[string]bool)
	for fromFormat := range cache.Conversions {
		formatSet[fromFormat] = true
		for toFormat := range cache.Conversions[fromFormat] {
			formatSet[toFormat] = true
		}
	}

	// Convert to sorted slice and filter out formats with no conversions
	allFormats := make([]string, 0, len(formatSet))
	for format := range formatSet {
		allFormats = append(allFormats, format)
	}
	sort.Strings(allFormats)

	// Build two sets: formats that can be sources and formats that can be targets
	sourceFormats := make(map[string]bool)
	targetFormats := make(map[string]bool)

	for fromFormat, toFormats := range cache.Conversions {
		hasOutgoing := false
		for toFormat, modes := range toFormats {
			for _, helpers := range modes {
				if len(helpers) > 0 {
					hasOutgoing = true
					targetFormats[toFormat] = true
				}
			}
		}
		if hasOutgoing {
			sourceFormats[fromFormat] = true
		}
	}

	// Keep formats that appear as BOTH source AND target
	// This removes columns with no incoming conversions and rows with no outgoing conversions
	formats := make([]string, 0)
	for _, format := range allFormats {
		if sourceFormats[format] && targetFormats[format] {
			formats = append(formats, format)
		}
	}

	if len(formats) == 0 {
		fmt.Println("No formats available")
		return
	}

	// Calculate column width based on longest format name
	maxFormatLen := 0
	for _, format := range formats {
		if len(format) > maxFormatLen {
			maxFormatLen = len(format)
		}
	}
	colWidth := maxFormatLen
	if colWidth < minColWidth {
		colWidth = minColWidth
	}

	// Print title
	fmt.Println("Format Conversion Matrix:")
	fmt.Println("(✓ = conversion supported)")
	fmt.Println()

	// Helper function to print separator row
	printSeparator := func() {
		fmt.Print(strings.Repeat("─", colWidth))
		fmt.Print("─┼")
		for i := range formats {
			fmt.Print(strings.Repeat("─", colWidth))
			if i < len(formats)-1 {
				fmt.Print("─┼")
			} else {
				fmt.Print("─")
			}
		}
		fmt.Println()
	}

	// Print header row
	fmt.Printf("%-*s │", colWidth, "FROM\\TO")
	for i, toFormat := range formats {
		fmt.Printf(" %-*s", colWidth-1, toFormat)
		if i < len(formats)-1 {
			fmt.Print(" │")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()

	// Print header separator
	printSeparator()

	// Print matrix rows
	for rowIdx, fromFormat := range formats {
		fmt.Printf("%-*s │", colWidth, fromFormat)
		for colIdx, toFormat := range formats {
			symbol := " "
			if fromFormat != toFormat {
				// Check if conversion exists (any mode)
				if toFormats, ok := cache.Conversions[fromFormat]; ok {
					if modes, ok := toFormats[toFormat]; ok && len(modes) > 0 {
						symbol = "✓"
					}
				}
			}
			fmt.Printf(" %-*s", colWidth-1, symbol)
			if colIdx < len(formats)-1 {
				fmt.Print(" │")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()

		// Print row separator (except after last row)
		if rowIdx < len(formats)-1 {
			printSeparator()
		}
	}
}
