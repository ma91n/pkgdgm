package main

import (
	"github.com/laqiiz/pkgdiagram/dependency"
	"github.com/laqiiz/pkgdiagram/directory"
	"github.com/laqiiz/pkgdiagram/umlgen"
	"log"
	"os"
	"path/filepath"
)

func main() {

	basePath := "/Users/mano/Go/src/github.com/laqiiz/gbilling-report"

	dirs, err := directory.NewWithIgnores(".git", ".idea").Do(basePath)
	if err != nil {
		log.Println(err)
	}

	for _, dir := range dirs {
		log.Printf("%+v\n", dir)
	}

	parser := dependency.New(filepath.Base(basePath))

	dependencies := make(dependency.Dependencies, 0)
	for _, v := range dirs {
		parsed := parser.Do(filepath.Join(basePath, v))
		dependencies = append(dependencies, parsed...)
	}
	for _, v := range dependencies {
		log.Printf("dependencies: %+v\n", v)
	}


	depends := dependencies.Pkgs()

	for _, v := range depends {
		log.Printf("depends: %+v\n", v)
	}

	file, err := os.Create("test.pu")
	if err != nil {
		log.Fatal(err)
	}

	if err := umlgen.New().Do(file, depends); err != nil {
		log.Fatal(err)
	}

	log.Println("finish")
}
