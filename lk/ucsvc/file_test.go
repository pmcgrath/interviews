package main

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	currentDir, _ := os.Getwd()

	testCases := []struct {
		Path     string
		Expected bool
	}{
		{"", false},
		{currentDir, false},
		{os.Args[0], true},
	}

	for _, testCase := range testCases {
		actual := fileExists(testCase.Path)
		if actual != testCase.Expected {
			t.Errorf("For %s expected %t, actual %t", testCase.Path, testCase.Expected, actual)
		}
	}
}
