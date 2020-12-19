package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNotEmptyString(t *testing.T) {
	args := []string{"", "path", "/my/path"}
	for _, arg := range args {
		NotEmptyString(arg)
	}
}

func TestFileExists(t *testing.T) {
	pathStr := "/testFile.txt"
	tempDir := os.TempDir()
	err := ioutil.WriteFile(tempDir+pathStr, []byte("Hello"), 07555)
	defer os.Remove(tempDir + pathStr)
	if err != nil {
		t.Fatalf("unexpected error while trying to test file exists %+v", err)
	}
	if FileNotExists(tempDir + pathStr) {
		t.Fatalf("File was suppoused to exist")
	}
}

func TestFileNotExists(t *testing.T) {
	pathStr := "/zero"
	tempDir := os.TempDir()
	if !FileNotExists(tempDir + pathStr) {
		t.Fatalf("File was suppoused to exist")
	}
}
