package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valpere/yakateka/internal/helper"
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

The cache file is used at runtime for fast helper lookup.`,
	RunE: runHelpers,
}

func init() {
	rootCmd.AddCommand(helpersCmd)
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

	return nil
}
