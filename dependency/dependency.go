package dependency

import (
	"go/parser"
	"go/token"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type pkkParser struct {
	Contains []string
}

func New(contains ...string) *pkkParser {
	return &pkkParser{
		Contains: contains,
	}
}

func (a *pkkParser) contains(s string) bool {
	for _, v := range a.Contains {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func (a *pkkParser) Do(path string) Dependencies {
	log.Println("parse dir is " + path)

	filter := func(info os.FileInfo) bool {
		if strings.HasSuffix(info.Name(), "_test.go") {
			return false
		}
		return true
	}

	fs := token.NewFileSet()
	f, err := parser.ParseDir(fs, path, filter, parser.Mode(parser.ImportsOnly))
	if err != nil {
		log.Fatal(err)
	}
	dependencies := make([]Dependency, 0)
	for _, v := range f {
		for path, astfile := range v.Files {
			relPath, err := filepath.Rel(filepath.Dir(path), path)
			if err != nil {
				log.Fatal(err)
			}

			imports := make([]string, 0, len(astfile.Imports))
			for _, imp := range astfile.Imports {
				path := imp.Path.Value
				if a.contains(path) {
					imports = append(imports, strings.Trim(path, "\""))
				}
			}

			dpn := Dependency{
				PackageName: astfile.Name.Name,
				FilePath:    relPath,
				Imports:     imports,
			}
			dependencies = append(dependencies, dpn)
		}
	}
	return dependencies
}

// Dependency represents go file import that exclude standard package.
type Dependency struct {
	PackageName string
	FilePath    string
	Imports     []string
	DependPkgs  []string
}

type Dependencies []Dependency

func (d Dependencies) searchBasePkg() string {
	minPath := ""
	minLen := math.MaxInt8
	for _, v := range d {
		for _, imp := range v.Imports {
			pkgLen := len(strings.Split(imp, "/"))
			if pkgLen < minLen {
				minLen = pkgLen
				minPath = imp
			}
		}
	}

	dir := filepath.Dir(minPath)
	return dir
}

func (d Dependencies) Pkgs() Dependencies {
	basePkg := d.searchBasePkg()
	log.Printf("basePkg=%s\n", basePkg)

	result := make([]Dependency, 0, len(d))

	for _, v := range d {
		dependPkgs := make([]string, 0, len(v.Imports))

		for _, imp := range v.Imports {
			dependPkg, err := filepath.Rel(basePkg, imp)
			if err != nil {
				// TODO temp impl
				panic(err)
			}
			dependPkgs = append(dependPkgs, dependPkg)
		}
		d := Dependency{
			PackageName: v.PackageName,
			FilePath:    v.FilePath,
			Imports:     v.Imports,
			DependPkgs:  dependPkgs,
		}
		result = append(result, d)
	}
	return result
}
