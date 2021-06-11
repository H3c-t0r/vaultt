// +build windows

package diagnose

import "io/fs"

// IsOwnedByRoot does nothing on windows
func IsOwnedByRoot(info fs.FileInfo) (bool, error) {
	return false, nil
}
