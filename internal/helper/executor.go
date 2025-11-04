package helper

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Executor executes helper commands
type Executor struct {
	timeout time.Duration
}

// NewExecutor creates a new helper executor
func NewExecutor(timeout time.Duration) *Executor {
	return &Executor{
		timeout: timeout,
	}
}

// Ping checks if helper is available
// Returns true if exit code == 0 and stdout == "pong"
func (e *Executor) Ping(ctx context.Context, helperPath string) bool {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, helperPath, "ping")
	output, err := cmd.Output()
	if err != nil {
		log.Debug().
			Err(err).
			Str("helper", helperPath).
			Msg("Helper ping failed")
		return false
	}

	response := strings.TrimSpace(string(output))
	if response != "pong" {
		log.Debug().
			Str("helper", helperPath).
			Str("response", response).
			Msg("Helper ping returned unexpected response")
		return false
	}

	return true
}

// GetInfo queries helper for its capabilities
// Returns HelperInfo on success, error if helper fails or returns invalid data
func (e *Executor) GetInfo(ctx context.Context, helperPath string) (*HelperInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, helperPath, "info")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Exit code > 0: error
	if err != nil {
		log.Warn().
			Err(err).
			Str("helper", helperPath).
			Str("stderr", stderr.String()).
			Msg("Helper info command failed")
		return nil, fmt.Errorf("helper returned error: %w - %s", err, stderr.String())
	}

	// Exit code == 0 but empty stdout: helper can't work now
	outputStr := strings.TrimSpace(stdout.String())
	if outputStr == "" {
		log.Warn().
			Str("helper", helperPath).
			Msg("Helper returned empty info (can't work now)")
		return nil, fmt.Errorf("helper returned empty info")
	}

	// Try to parse as YAML
	var info HelperInfo
	if err := yaml.Unmarshal([]byte(outputStr), &info); err != nil {
		log.Error().
			Err(err).
			Str("helper", helperPath).
			Str("output", outputStr).
			Msg("Helper returned invalid YAML")
		return nil, fmt.Errorf("invalid YAML from helper: %w", err)
	}

	// Validate that we got some capabilities
	if len(info.Capabilities) == 0 {
		log.Warn().
			Str("helper", helperPath).
			Msg("Helper returned empty capabilities")
		return nil, fmt.Errorf("helper has no capabilities")
	}

	log.Debug().
		Str("helper", helperPath).
		Str("name", info.Name).
		Int("formats", len(info.Capabilities)).
		Msg("Successfully loaded helper info")

	return &info, nil
}

// Convert executes a conversion using the helper
func (e *Executor) Convert(ctx context.Context, helperPath string, mode ConversionMode, fromFormat, fromFile, toFormat, toFile string) error {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, helperPath, "convert", string(mode), fromFormat, fromFile, toFormat, toFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Error().
			Err(err).
			Str("helper", helperPath).
			Str("mode", string(mode)).
			Str("from", fromFormat).
			Str("to", toFormat).
			Str("stderr", stderr.String()).
			Msg("Helper conversion failed")
		return fmt.Errorf("conversion failed: %w - %s", err, stderr.String())
	}

	log.Info().
		Str("helper", helperPath).
		Str("mode", string(mode)).
		Str("from", fromFormat).
		Str("to", toFormat).
		Str("input", fromFile).
		Str("output", toFile).
		Msg("Successfully converted with helper")

	return nil
}
