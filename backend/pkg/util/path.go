package util

import (
	"os"
	"path/filepath"
)

// GetCoursesDirectory figures out where course files are stored
func GetCoursesDirectory() string {
	// check container env var first
	coursesDir := os.Getenv("INTERNAL_COURSES_DIR")
	if coursesDir != "" {
		return coursesDir
	}

	// fallback to local dev env var
	coursesDir = os.Getenv("COURSES_BASE_DIR")
	if coursesDir == "" {
		// last resort - current directory
		coursesDir = "."
	}

	return coursesDir
}

// EnsureDirectoryExists creates directory if it doesn't exist
func EnsureDirectoryExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// try to create it
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return false
		}
	}
	return true
}

// ResolveCourseFilePath converts relative path to absolute path for course files
func ResolveCourseFilePath(relativePath string) string {
	baseDir := GetCoursesDirectory()
	return filepath.Join(baseDir, relativePath)
}
