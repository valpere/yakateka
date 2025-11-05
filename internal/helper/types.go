package helper

import (
	"github.com/valpere/yakateka/internal"
)

// ConversionMode represents the mode of conversion
type ConversionMode string

const (
	ModeNormal  ConversionMode = "normal"  // Balanced speed and quality (mandatory)
	ModeFast    ConversionMode = "fast"    // Speed over quality (optional)
	ModeQuality ConversionMode = "quality" // Quality over speed (optional)
)

// ModeMetrics contains performance metrics for a conversion mode
type ModeMetrics struct {
	Speed   float64 `yaml:"speed"`   // > 0 means supported (higher = faster)
	Quality float64 `yaml:"quality"` // > 0 means supported (higher = better)
}

// IsSupported returns true if this mode is supported (both metrics > 0)
func (m ModeMetrics) IsSupported() bool {
	return m.Speed > 0 && m.Quality > 0
}

// FormatPair represents conversion from one format to another
type FormatPair struct {
	From internal.DocumentFormat
	To   internal.DocumentFormat
}

// ModeCapabilities contains all modes for a format conversion
// Wraps the modes in a "modes" key to match helper output structure
type ModeCapabilities struct {
	Modes ModesStruct `yaml:"modes"`
}

// ModesStruct contains the actual mode definitions
type ModesStruct struct {
	Normal  ModeMetrics `yaml:"normal"`           // Mandatory
	Fast    ModeMetrics `yaml:"fast,omitempty"`   // Optional
	Quality ModeMetrics `yaml:"quality,omitempty"` // Optional
}

// GetMode returns metrics for specified mode, with fallback to normal
func (mc *ModeCapabilities) GetMode(mode ConversionMode) ModeMetrics {
	switch mode {
	case ModeFast:
		if mc.Modes.Fast.IsSupported() {
			return mc.Modes.Fast
		}
	case ModeQuality:
		if mc.Modes.Quality.IsSupported() {
			return mc.Modes.Quality
		}
	}
	// Fallback to normal
	return mc.Modes.Normal
}

// HelperInfo contains information returned by helper.sh info
type HelperInfo struct {
	Name         string                                       `yaml:"name"`
	Version      string                                       `yaml:"version,omitempty"`
	Description  string                                       `yaml:"description,omitempty"`
	Capabilities map[string]map[string]ModeCapabilities      `yaml:"capabilities"` // from_format -> to_format -> modes
}

// HelperConfig represents a helper's configuration
type HelperConfig struct {
	Path      string  // Path to helper script
	Weight    float64 // Weight for prioritization (0.0 to 1.0)
	Available bool    // Whether helper passed ping check
}

// HelperEntry represents a helper in the registry
type HelperEntry struct {
	Config HelperConfig
	Info   *HelperInfo
}

// Registry contains all registered helpers
type Registry struct {
	Helpers map[string]*HelperEntry // path -> entry
}

// CacheEntry represents a helper entry in helpers.yaml cache
type CacheEntry struct {
	Helper string  `yaml:"helper"`
	Weight float64 `yaml:"weight"`
	// Future: can add speed/quality metrics here
}

// HelperCache represents the helpers.yaml file structure
type HelperCache struct {
	// Structure: from_format -> to_format -> mode -> []CacheEntry
	Conversions map[string]map[string]map[string][]CacheEntry `yaml:"conversions"`
}
