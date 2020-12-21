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
	err := createFileHelper(tempDir+pathStr, "Hello ")
	defer os.Remove(tempDir + pathStr)
	if err != nil {
		t.Fatalf("unexpected error while trying to test file exists %+v", err)
	}
	if FileNotExists(tempDir + pathStr) {
		t.Fatalf("File was suppoused to exist")
	}
}

func createFileHelper(path, content string) (err error) {
	err = ioutil.WriteFile(path, []byte(content), 07555)
	return
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
	scenarios := [...]struct {
		data       string
		def        Definition
		expectsErr bool
	}{
		{
			data:       "definitions:\n  foo:\n    path: mypath\n    template: fakePath\n",
			def:        Definition{Definitions: make(map[string]Tool, 1)},
			expectsErr: false,
		},
		{
			data:       "xasdfsadflkn",
			def:        Definition{Definitions: make(map[string]Tool, 0)},
			expectsErr: true,
		},
	}
	for _, scenario := range scenarios {
		_, err := YMLStringToDefinitions(scenario.data)
		if !IsNotNil(err) == scenario.expectsErr {
			t.Fatalf("expected %v got %v", scenario.expectsErr, err)
		}
	}
}
func TestFailIfError(t *testing.T) {
	FailIfError("must fail if a not nil error is passed", nil)
}

func TestRemoveIfExists(t *testing.T) {
	scenarios := [...]struct {
		def            Tool
		mustCreateFile bool
		errorExpected  bool
	}{
		{
			def: Tool{
				Template: "testingTemplate",
				Path:     "/p.txt",
			},
			mustCreateFile: true,
			errorExpected:  false,
		},
		{
			def: Tool{
				Path:     "mustFail",
				Template: "/p.txt",
			},
			mustCreateFile: false,
			errorExpected:  true,
		},
	}
	for _, scenario := range scenarios {
		if scenario.mustCreateFile {
			_ = createFileHelper(os.TempDir()+scenario.def.Path, "Hello")
		}
		def := Tool{
			Path:     os.TempDir() + scenario.def.Path,
			Template: scenario.def.Template,
		}
		err := RemovesFileIfExists(def)
		if !IsNotNil(err) == scenario.errorExpected {
			t.Fatalf("file should be deleted without any issue")
		}
	}
}

func TestCreateSymlink(t *testing.T) {
	tDir := "/fake"
	os.Mkdir(os.TempDir()+tDir, os.ModePerm)
	createFileHelper(os.TempDir()+tDir+"/text.tempalte", "Hello")
	def := Tool{
		Path:     os.TempDir() + "/fake/text.symlink",
		Template: "/text.template",
	}
	err := CreateSymlinkFromDef(os.TempDir(), def)
	if IsNotNil(err) {
		t.Fatalf("error while trying to create symlink")
	}
	oldData, _ := ioutil.ReadFile(os.TempDir() + "/fake/text.symlink")
	symData, _ := ioutil.ReadFile(os.TempDir() + "/fake/text.template")
	if string(oldData) != string(symData) {
		t.Fatalf("data expeceted %s got %s", oldData, symData)
	}
	defer os.RemoveAll("/tmp/fake/")
}

func TestFindAndReadDefinitions(t *testing.T) {
	tDir := "/fakedef"
	expectedText := "Hello world"
	os.Mkdir(os.TempDir()+tDir, os.ModePerm)
	createFileHelper(os.TempDir()+tDir+"/definitions.yml", expectedText)
	defString := FindAndReadDefinitions(os.TempDir() + tDir + "/")
	if expectedText != defString {
		t.Fatalf("expected %s got %s ", expectedText, defString)
	}
	defer os.RemoveAll(os.TempDir() + tDir)
}
