package main

import "os"

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
