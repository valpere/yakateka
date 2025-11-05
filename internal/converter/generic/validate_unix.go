//go:build unix || linux || darwin

package generic

import (
	"os"

	"github.com/rs/zerolog/log"
)

// checkExecutable checks if a file has execute permissions (Unix-specific)
func checkExecutable(info os.FileInfo, binaryPath string) error {
	// On Unix-like systems, check if file is executable
	if info.Mode()&0111 == 0 {
		log.Warn().
			Str("binary", binaryPath).
			Msg("Binary may not be executable (no execute permission)")
	}
	return nil
}
