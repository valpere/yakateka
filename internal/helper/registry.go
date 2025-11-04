package helper

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valpere/yakateka/internal"
	"gopkg.in/yaml.v3"
)

// NewRegistry creates a new helper registry
func NewRegistry() *Registry {
	return &Registry{
		Helpers: make(map[string]*HelperEntry),
	}
}

// Register adds a helper to the registry
func (r *Registry) Register(path string, weight float64) {
	r.Helpers[path] = &HelperEntry{
		Config: HelperConfig{
			Path:      path,
			Weight:    weight,
			Available: false, // Will be set by ping
		},
		Info: nil, // Will be populated by GetInfo
	}
}

// Initialize pings all helpers and loads their info
func (r *Registry) Initialize(ctx context.Context, executor *Executor) error {
	for path, entry := range r.Helpers {
		// Ping helper
		entry.Config.Available = executor.Ping(ctx, path)
		if !entry.Config.Available {
			log.Warn().Str("helper", path).Msg("Helper failed ping check")
			continue
		}

		// Get info
		info, err := executor.GetInfo(ctx, path)
		if err != nil {
			log.Warn().
				Err(err).
				Str("helper", path).
				Msg("Failed to get helper info")
			entry.Config.Available = false
			continue
		}

		entry.Info = info
		log.Info().
			Str("helper", path).
			Str("name", info.Name).
			Float64("weight", entry.Config.Weight).
			Msg("Registered helper")
	}

	return nil
}

// GenerateCache creates helpers.yaml cache file
func (r *Registry) GenerateCache() (*HelperCache, error) {
	cache := &HelperCache{
		Conversions: make(map[string]map[string]map[string][]CacheEntry),
	}

	// Collect all conversions from all helpers
	type conversionKey struct {
		from string
		to   string
		mode string
	}

	conversions := make(map[conversionKey][]CacheEntry)

	// Iterate through all helpers
	for path, entry := range r.Helpers {
		if !entry.Config.Available || entry.Info == nil {
			continue
		}

		// Iterate through capabilities
		for fromFormat, toFormats := range entry.Info.Capabilities {
			for toFormat, modes := range toFormats {
				// Check each mode
				for _, modeName := range []ConversionMode{ModeNormal, ModeFast, ModeQuality} {
					var metrics ModeMetrics
					switch modeName {
					case ModeNormal:
						metrics = modes.Modes.Normal
					case ModeFast:
						metrics = modes.Modes.Fast
					case ModeQuality:
						metrics = modes.Modes.Quality
					}

					if metrics.IsSupported() {
						key := conversionKey{
							from: fromFormat,
							to:   toFormat,
							mode: string(modeName),
						}

						conversions[key] = append(conversions[key], CacheEntry{
							Helper: path,
							Weight: entry.Config.Weight,
						})
					}
				}
			}
		}
	}

	// Sort each conversion list by weight (descending) then by name
	for key, entries := range conversions {
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Weight != entries[j].Weight {
				return entries[i].Weight > entries[j].Weight // Descending
			}
			return entries[i].Helper < entries[j].Helper // Alphabetical
		})

		// Add to cache structure
		if cache.Conversions[key.from] == nil {
			cache.Conversions[key.from] = make(map[string]map[string][]CacheEntry)
		}
		if cache.Conversions[key.from][key.to] == nil {
			cache.Conversions[key.from][key.to] = make(map[string][]CacheEntry)
		}
		cache.Conversions[key.from][key.to][key.mode] = entries
	}

	return cache, nil
}

// SaveCache saves helpers.yaml to disk
func (cache *HelperCache) SaveCache(path string) error {
	data, err := yaml.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	log.Info().Str("path", path).Msg("Saved helper cache")
	return nil
}

// LoadCache loads helpers.yaml from disk
func LoadCache(path string) (*HelperCache, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var cache HelperCache
	if err := yaml.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	return &cache, nil
}

// FindHelpers returns ordered list of helpers for a conversion
// Falls back from requested mode to normal if needed
func (cache *HelperCache) FindHelpers(from, to internal.DocumentFormat, mode ConversionMode) []CacheEntry {
	fromStr := string(from)
	toStr := string(to)
	modeStr := string(mode)

	// Try requested mode first
	if helpers, ok := cache.Conversions[fromStr][toStr][modeStr]; ok && len(helpers) > 0 {
		return helpers
	}

	// Fallback to normal mode
	if mode != ModeNormal {
		if helpers, ok := cache.Conversions[fromStr][toStr][string(ModeNormal)]; ok {
			log.Debug().
				Str("from", fromStr).
				Str("to", toStr).
				Str("requested", modeStr).
				Msg("No helpers for requested mode, falling back to normal")
			return helpers
		}
	}

	return nil
}

// MarkHelperFailed marks a helper as unavailable for a specific conversion
// This prevents retrying the same helper for the same conversion pair
func (cache *HelperCache) MarkHelperFailed(from, to internal.DocumentFormat, helperPath string) {
	fromStr := string(from)
	toStr := string(to)

	if cache.Conversions[fromStr] == nil {
		return
	}
	if cache.Conversions[fromStr][toStr] == nil {
		return
	}

	// Remove helper from all modes for this conversion
	for mode, helpers := range cache.Conversions[fromStr][toStr] {
		filtered := make([]CacheEntry, 0, len(helpers))
		for _, h := range helpers {
			if h.Helper != helperPath {
				filtered = append(filtered, h)
			}
		}
		cache.Conversions[fromStr][toStr][mode] = filtered
	}

	log.Debug().
		Str("from", fromStr).
		Str("to", toStr).
		Str("helper", helperPath).
		Msg("Marked helper as failed for this conversion")
}

// GetTimeout returns default timeout for helper operations
func GetTimeout() time.Duration {
	return 10 * time.Second // Default timeout for info/ping
}
