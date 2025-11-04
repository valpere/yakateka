package config

// ConverterConfig defines the complete converter configuration
type ConverterConfig struct {
	Profiles   map[string]ProfileConfig `mapstructure:"converter_profiles" yaml:"converter_profiles"`
	Converters map[string]ToolConfig    `mapstructure:"converters" yaml:"converters"`
}

// ProfileConfig defines a reusable command pattern
type ProfileConfig struct {
	CommandTemplate string `mapstructure:"command_template" yaml:"command_template"`
	PostProcess     string `mapstructure:"post_process" yaml:"post_process"`
}

// ToolConfig defines a converter tool configuration
type ToolConfig struct {
	Binary              string                        `mapstructure:"binary" yaml:"binary"`
	Profile             string                        `mapstructure:"profile" yaml:"profile"`
	CommandTemplate     string                        `mapstructure:"command_template" yaml:"command_template"` // Override profile
	Timeout             int                           `mapstructure:"timeout" yaml:"timeout"`
	Formats             FormatConfig                  `mapstructure:"formats" yaml:"formats"`
	FormatMapping       map[string]string             `mapstructure:"format_mapping" yaml:"format_mapping"`
	ConversionOverrides map[string]ConversionOverride `mapstructure:"conversion_overrides" yaml:"conversion_overrides"`
}

// FormatConfig defines supported input and output formats
type FormatConfig struct {
	Input  []string `mapstructure:"input" yaml:"input"`
	Output []string `mapstructure:"output" yaml:"output"`
}

// ConversionOverride defines format-specific conversion settings
type ConversionOverride struct {
	ExtraArgs    string            `mapstructure:"extra_args" yaml:"extra_args"`
	OutputFormat string            `mapstructure:"output_format" yaml:"output_format"`
	Quality      map[string]string `mapstructure:"quality" yaml:"quality"`
}

// GetCommandTemplate returns the command template for this tool
// Uses tool-specific template if set, otherwise falls back to profile
func (tc *ToolConfig) GetCommandTemplate(profiles map[string]ProfileConfig) string {
	if tc.CommandTemplate != "" {
		return tc.CommandTemplate
	}
	if tc.Profile != "" {
		if profile, ok := profiles[tc.Profile]; ok {
			return profile.CommandTemplate
		}
	}
	return ""
}

// GetPostProcess returns the post-processing action for this tool
func (tc *ToolConfig) GetPostProcess(profiles map[string]ProfileConfig) string {
	if tc.Profile != "" {
		if profile, ok := profiles[tc.Profile]; ok {
			return profile.PostProcess
		}
	}
	return ""
}

// GetConversionOverride returns override for specific conversion
// Supports wildcards: "*->pdf", "md->*", "*->*"
func (tc *ToolConfig) GetConversionOverride(inputFormat, outputFormat string) *ConversionOverride {
	key := inputFormat + "->" + outputFormat

	// Try exact match first
	if override, ok := tc.ConversionOverrides[key]; ok {
		return &override
	}

	// Try input wildcard: *->output
	if override, ok := tc.ConversionOverrides["*->"+outputFormat]; ok {
		return &override
	}

	// Try output wildcard: input->*
	if override, ok := tc.ConversionOverrides[inputFormat+"->*"]; ok {
		return &override
	}

	// Try full wildcard: *->*
	if override, ok := tc.ConversionOverrides["*->*"]; ok {
		return &override
	}

	return nil
}

// MapFormat maps internal format name to tool-specific format name
func (tc *ToolConfig) MapFormat(format string) string {
	if mapped, ok := tc.FormatMapping[format]; ok {
		return mapped
	}
	return format
}

// SupportsConversion checks if this tool supports a given conversion
func (tc *ToolConfig) SupportsConversion(inputFormat, outputFormat string) bool {
	hasInput := false
	hasOutput := false

	for _, f := range tc.Formats.Input {
		if f == inputFormat {
			hasInput = true
			break
		}
	}

	for _, f := range tc.Formats.Output {
		if f == outputFormat {
			hasOutput = true
			break
		}
	}

	return hasInput && hasOutput
}
