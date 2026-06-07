package pathutil

import (
	"os"
	"path/filepath"
)

func ResolveExistingRelativeFile(relativePath string) (string, bool) {
	if filepath.IsAbs(relativePath) {
		if _, err := os.Stat(relativePath); err == nil {
			return relativePath, true
		}
		return "", false
	}
	if exePath, err := os.Executable(); err == nil {
		if path, ok := FindFileFromBase(filepath.Dir(exePath), relativePath); ok {
			return path, true
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		if path, ok := FindFileFromBase(cwd, relativePath); ok {
			return path, true
		}
	}
	if abs, err := filepath.Abs(relativePath); err == nil {
		if _, statErr := os.Stat(abs); statErr == nil {
			return abs, true
		}
	}
	return "", false
}

func ResolveFilePath(relativePath string) string {
	if path, ok := ResolveExistingRelativeFile(relativePath); ok {
		return path
	}
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	if abs, err := filepath.Abs(relativePath); err == nil {
		return abs
	}
	return relativePath
}

func FindFileFromBase(baseDir, relativePath string) (string, bool) {
	dir := filepath.Clean(baseDir)
	for {
		candidate := filepath.Join(dir, relativePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}
