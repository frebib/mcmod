package util

import (
	"errors"
	"path"
)

var (
	ErrConflictingPaths = errors.New("conflicting output path and directory provided")
)

func StringInSlice(haystack []string, needle string) bool {
	for _, str := range haystack {
		if str == needle {
			return true
		}
	}
	return false
}

func CalcFilePath(name, ovFile, ovDir string) (string, error) {
	if ovFile != "" {
		// Error if given two file paths
		ovFileDir, _ := path.Split(ovFile)
		// Two directory paths are provided, we can't know which one to pick
		if ovFileDir != "" && ovDir != "" {
			return "", ErrConflictingPaths
		}
		if ovFile == "-" {
			ovFile = "/dev/stdout"
		}
		name = ovFile
	}
	return path.Join(ovDir, name), nil
}
