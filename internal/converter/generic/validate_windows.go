//go:build windows

package generic

import (
	"os"
)

// checkExecutable checks if a file is executable (Windows version - no-op)
// On Windows, file extensions (.exe, .bat, .cmd) determine executability
func checkExecutable(info os.FileInfo, binaryPath string) error {
	// Windows doesn't use Unix permission bits
	// Executability is determined by file extension
	return nil
}
