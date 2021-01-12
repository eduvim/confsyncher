package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// DefinitionFilename is the name of the file expected to contain the definition yml
const DefinitionFilename = "definitions.yml"

// Definition wraps the complete map of definitions
type Definition struct {
	Definitions map[string]Tool
}

// Tool maps the path where you need the template file to be and the path for your template
type Tool struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

func main() {
	execute()
}

func execute() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("you need to provide the path where your definitions.yml directory")
	}
	data := FindAndReadDefinitions(args[1])
	usersDefs, err := YMLStringToDefinitions(data)
	FailIfError("the definitions file cannot be correctly parsed ", err)
	for toolName, defs := range usersDefs.Definitions {
		log.Printf("writing ...  %s \n", toolName)
		err := RemovesFileIfExists(defs)
		FailIfError("the definitions file cannot be correctly parsed ", err)
		err = CreateSymlinkFromDef(args[1], defs)
		FailIfError("the definitions file cannot be correctly parsed ", err)
		log.Printf("tool configured ...  %s \n", toolName)
	}
}

// IsNotNil returns true if the error sended is different to nil
func IsNotNil(value error) bool {
	return value != nil
}

// NotEmptyString trims the string and returns true if the arg is not empty
func NotEmptyString(arg string) bool {
	return len(strings.TrimSpace(arg)) > 0
}

// FileNotExists Returns true if the file doesn't exists
func FileNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// FindAndReadDefinitions looks into the dir passed the definitions.yml file and read it as string
func FindAndReadDefinitions(path string) (data string) {
	files, err := ioutil.ReadDir(path)
	if IsNotNil(err) {
		log.Fatalf("cannot open the declared dir %s %+v", path, err)
	}
	for _, file := range files {
		if file.Name() == DefinitionFilename {
			bytes, err := ioutil.ReadFile(path + DefinitionFilename)
			data = string(bytes)
			if IsNotNil(err) {
				log.Fatalf("cannot open declaration file %+v \n", err)
			}
			return data
		}
	}
	log.Fatalf("declaration file not found")
	return
}

// YMLStringToDefinitions binds the string value to a Definitions struct
func YMLStringToDefinitions(data string) (definition Definition, err error) {
	err = yaml.Unmarshal([]byte(data), &definition)
	return
}

// CreateSymlinkFromDef create a symlink to the old path and the definided path
func CreateSymlinkFromDef(path string, defs Tool) (err error) {
	err = os.Symlink(path+defs.Template, defs.Path)
	return
}

// RemovesFileIfExists removes a file on the destination path if it exists
func RemovesFileIfExists(defs Tool) (err error) {
	if !FileNotExists(defs.Path) {
		err = os.Remove(defs.Path)
	} else {
		err = errors.New("the file doesn't exist")
	}
	return
}

// FailIfError prints the error base on the text message , it will exit the program
func FailIfError(text string, err error) {
	if IsNotNil(err) {
		log.Fatalf(text, err)
	}
}
