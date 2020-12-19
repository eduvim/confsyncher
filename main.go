package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Definition struct {
	Definitions map[string]Tool
}
type Tool struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

func main() {
	args := os.Args
	files, err := ioutil.ReadDir(args[1])
	if IsNotNil(err) {
		log.Fatalf("cannot open the declared dir %s %+v", args[1], err)
	}

	data := ""
	for _, file := range files {
		if file.Name() == "definitions.yml" {
			bytes, err := ioutil.ReadFile(os.Args[1] + "/definitions.yml")
			data = string(bytes)
			if IsNotNil(err) {
				log.Fatalf("cannot open declaration file %+v \n", err)
			}
			log.Printf("Definitions file founded")
			break
		}
	}
	t := Definition{}
	err = yaml.Unmarshal([]byte(data), &t)
	if IsNotNil(err) {
		log.Fatalf("error: %v", err)
	}
	for toolName, defs := range t.Definitions {
		log.Printf("writing ...  %s \n", toolName)
		if !FileNotExists(defs.Path) {
			err = os.Remove(defs.Path)
			if IsNotNil(err) {
				log.Fatalf(" cannot remove path file due an error %v \n", err)
			}
		}
		log.Printf("path: %s template: %s", defs.Path, defs.Template)
		err = os.Symlink(args[1]+defs.Template, defs.Path)
		if IsNotNil(err) {
			log.Printf("cannot symlink due an error %v \n", err)
			return
		}
		log.Printf("tool configured ...  %s \n", toolName)
	}
}

func IsNotNil(value error) bool {
	return value != nil
}

func NotEmptyString(arg string) bool {
	return len(strings.TrimSpace(arg)) > 0
}

func FileNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
