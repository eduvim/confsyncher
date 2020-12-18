package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Definition struct {
	Definitions map[string]Tool
}
type Tool struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

func main() {
	// validate file
	args := os.Args
	fmt.Println(args[0])
	files, err := ioutil.ReadDir(args[1])
	if IsNotNil(err){
		log.Fatalf("cannot open the declared dir %s %+v" , args[1] , err)
	}
	data := ""
	for _ , file := range files{
		if file.Name() == "definitions.yml"{
			bytes , err := ioutil.ReadFile(os.Args[1] + "/definitions.yml")
			data = string(bytes)
			if IsNotNil(err) {
				log.Fatalf("cannot open declaration file %+v \n"  , err)
			}
			fmt.Println("Definitions file founded")
			break
		}
	}

	t := Definition{}
	err = yaml.Unmarshal([]byte(data), &t)
	if IsNotNil(err) {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)
	for toolName, defs := range t.Definitions{
		log.Printf("writing ...  %s \n" , toolName)

		//err = os.Remove(defs.Path)
		//if IsNotNil(err) {
		//	fmt.Printf(" cannot remove path file due an error %v \n", err)
		//	return
		//}

		log.Printf("path: %s template: %s" , defs.Path , defs.Template)
		err = os.Symlink( args[1] + defs.Template , defs.Path)
		if IsNotNil(err) {
			fmt.Printf("cannot symlink due an error %v \n", err)
			return
		}
		fmt.Printf("tool configured ...  %s \n" , toolName)
	}
}

func IsNotNil(value error) bool{
	return value != nil
}

func NotEmptyString(arg string) bool {
	return len(strings.TrimSpace(arg)) > 0
}
