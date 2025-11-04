package cmd

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that root command can be initialized
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	// Test command name
	if rootCmd.Use != "yakateka" {
		t.Errorf("Expected command name 'yakateka', got '%s'", rootCmd.Use)
	}

	// Test that persistent flags are registered
	flags := rootCmd.PersistentFlags()

	if !flags.HasFlags() {
		t.Error("Expected persistent flags to be registered")
	}

	// Test config flag
	configFlag := flags.Lookup("config")
	if configFlag == nil {
		t.Error("Expected 'config' flag to be registered")
	}

	// Test log-level flag
	logLevelFlag := flags.Lookup("log-level")
	if logLevelFlag == nil {
		t.Error("Expected 'log-level' flag to be registered")
	}

	// Test log-format flag
	logFormatFlag := flags.Lookup("log-format")
	if logFormatFlag == nil {
		t.Error("Expected 'log-format' flag to be registered")
	}

	// Test verbose flag
	verboseFlag := flags.Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected 'verbose' flag to be registered")
	}
}

func TestSetDefaults(t *testing.T) {
	// Initialize config
	initConfig()

	// This function sets defaults via viper
	// Just verify it doesn't panic
	setDefaults()

	// Test would be more comprehensive with actual viper value checks
	// but that requires more setup for testing
}
