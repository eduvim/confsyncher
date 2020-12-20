package main

import (
	"errors"
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

func TestIsNotNil(t *testing.T) {
	cases := [...]struct {
		err      error
		expected bool
	}{
		{
			err:      nil,
			expected: false,
		},
		{
			err:      errors.New("foo error"),
			expected: true,
		},
	}
	for _, val := range cases {
		if val.expected != IsNotNil(val.err) {
			t.Fatalf("want %v got %v ", val.expected, val.err)
		}
	}
}

func TestYMLStringToDefinitions(t *testing.T) {
	scenarios := []struct {
		data string
		def  Definition
	}{
		{
			data: "definitions:\n  foo:\n    path: mypath\n    template: fakePath\n",
			def:  Definition{Definitions: make(map[string]Tool, 1)},
		},
	}
	for _, scenario := range scenarios {
		YMLStringToDefinitions(scenario.data)
	}
}
