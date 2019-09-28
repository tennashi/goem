package shellpath

import (
	"os"
	"path/filepath"
	"strings"
)

func Resolve(path string) string {
	path = os.ExpandEnv(path)
	path = filepath.Clean(path)

	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}

	return path
}
