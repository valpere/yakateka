package helper

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// LoadAndPing loads helpers.yaml cache and pings all helpers
// Returns nil if cache doesn't exist or is empty
func LoadAndPing(ctx context.Context) (*HelperConverter, error) {
	// Get cache file path from config
	cacheFile := viper.GetString("helpers.cache_file")
	if cacheFile == "" {
		cacheFile = "helpers.yaml"
	}

	// Check if cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		log.Debug().
			Str("cache", cacheFile).
			Msg("Helper cache file not found, skipping helper system")
		return nil, nil
	}

	// Load cache
	cache, err := LoadCache(cacheFile)
	if err != nil {
		log.Warn().
			Err(err).
			Str("cache", cacheFile).
			Msg("Failed to load helper cache")
		return nil, err
	}

	if len(cache.Conversions) == 0 {
		log.Debug().Msg("Helper cache is empty, skipping helper system")
		return nil, nil
	}

	log.Info().
		Str("cache", cacheFile).
		Msg("Loaded helper cache")

	// Ping all unique helpers
	helperPaths := make(map[string]bool)
	for _, toFormats := range cache.Conversions {
		for _, modes := range toFormats {
			for _, helpers := range modes {
				for _, h := range helpers {
					helperPaths[h.Helper] = false // false = not yet pinged
				}
			}
		}
	}

	executor := NewExecutor(GetTimeout())
	successCount := 0
	failCount := 0

	for helperPath := range helperPaths {
		if executor.Ping(ctx, helperPath) {
			helperPaths[helperPath] = true
			successCount++
			log.Debug().
				Str("helper", helperPath).
				Msg("Helper ping successful")
		} else {
			failCount++
			log.Warn().
				Str("helper", helperPath).
				Msg("Helper ping failed, will be excluded")

			// Mark helper as globally failed (more efficient than per-conversion)
			cache.MarkHelperGloballyFailed(helperPath)
		}
	}

	log.Info().
		Int("available", successCount).
		Int("unavailable", failCount).
		Msg("Helper ping check complete")

	// If all helpers failed, return nil
	if successCount == 0 {
		log.Warn().Msg("No helpers available, helper system disabled")
		return nil, nil
	}

	// Create and return converter
	converter := NewHelperConverter(cache, executor)
	return converter, nil
}
