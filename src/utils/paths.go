package utils

import (
	"path/filepath"
	"strings"
)

func FindRootDirectoryOfAllPaths(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	rootDir := filepath.Dir(paths[0])

	for _, path := range paths[1:] {
		for !strings.HasPrefix(path, rootDir) {
			rootDir = filepath.Dir(rootDir)
		}
	}

	return rootDir
}
