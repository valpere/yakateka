package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load loads converter configuration from viper
func Load() (*ConverterConfig, error) {
	var cfg ConverterConfig

	// Try to unmarshal converter configuration
	if err := viper.UnmarshalKey("converter_profiles", &cfg.Profiles); err != nil {
		return nil, fmt.Errorf("failed to load converter profiles: %w", err)
	}

	if err := viper.UnmarshalKey("converters", &cfg.Converters); err != nil {
		return nil, fmt.Errorf("failed to load converters: %w", err)
	}

	// Validate configuration
	if err := Validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid converter configuration: %w", err)
	}

	return &cfg, nil
}

// Validate validates the converter configuration
func Validate(cfg *ConverterConfig) error {
	// Validate profiles
	for name, profile := range cfg.Profiles {
		if profile.CommandTemplate == "" {
			return fmt.Errorf("profile %s has empty command_template", name)
		}
	}

	// Validate converters
	for name, tool := range cfg.Converters {
		// Must have either profile or command_template
		if tool.Profile == "" && tool.CommandTemplate == "" {
			return fmt.Errorf("converter %s must have either profile or command_template", name)
		}

		// If using profile, it must exist
		if tool.Profile != "" {
			if _, ok := cfg.Profiles[tool.Profile]; !ok {
				return fmt.Errorf("converter %s references unknown profile: %s", name, tool.Profile)
			}
		}

		// Must have binary path
		if tool.Binary == "" {
			return fmt.Errorf("converter %s has empty binary", name)
		}

		// Must have at least one input and output format
		if len(tool.Formats.Input) == 0 {
			return fmt.Errorf("converter %s has no input formats", name)
		}
		if len(tool.Formats.Output) == 0 {
			return fmt.Errorf("converter %s has no output formats", name)
		}
	}

	return nil
}
