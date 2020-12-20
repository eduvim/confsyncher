package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

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
	usersDefs := YMLStringToDefinitions(data)
	for toolName, defs := range usersDefs.Definitions {
		log.Printf("writing ...  %s \n", toolName)
		RemovesFileIfExists(defs)
		CreateSymlinkFromDef(args[1], defs)
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
		if file.Name() == "definitions.yml" {
			bytes, err := ioutil.ReadFile(path + "/definitions.yml")
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
func YMLStringToDefinitions(data string) (definition Definition) {
	err := yaml.Unmarshal([]byte(data), &definition)
	if IsNotNil(err) {
		log.Fatalf("error: %v", err)
	}
	return
}

// CreateSymlinkFromDef create a symlink to the old path and the definided path
func CreateSymlinkFromDef(path string, defs Tool) {
	log.Printf("path: %s template: %s", defs.Path, defs.Template)
	err := os.Symlink(path+defs.Template, defs.Path)
	if IsNotNil(err) {
		log.Printf("cannot symlink due an error %v \n", err)
		return
	}

}

// RemovesFileIfExists removes a file on the destination path if it exists
func RemovesFileIfExists(defs Tool) {
	if !FileNotExists(defs.Path) {
		err := os.Remove(defs.Path)
		if IsNotNil(err) {
			log.Fatalf(" cannot remove path file due an error %v \n", err)
		}
	}
}
