package backend

import (
	"errors"
	"os"
)

type PathDir struct {
	Home string
}

func (*PathDir) GetHomePath() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// Prefer standard environment variable USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}

	return "", errors.New("未找到主目录")
}

func NewPathDir() *PathDir {
	return &PathDir{}
}
