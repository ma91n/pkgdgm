package umlgen

import (
	"fmt"
	"io"
	"strings"

	"github.com/laqiiz/pkgdiagram/dependency"
)

type generator struct {
}

func New() *generator {
	return &generator{}
}

func (g *generator) Do(w io.Writer, dependencies []dependency.Dependency) error {

	var b strings.Builder
	b.WriteString("@startuml\n")
	b.WriteString("title package-diagram\n\n")

	// distinct package name
	pkgSet := map[string]bool{}
	for _, v := range dependencies {
		pkgSet[v.PackageName] = true
	}

	for k := range pkgSet {
		_, err := fmt.Fprintf(&b, "package %s{}\n", k)
		if err != nil {
			return err
		}
	}
	b.WriteString("\n")

	// link
	// distinct dependency relations

	packageDependencyMap := map[string][]dependency.Dependency{}
	for _, v := range dependencies {
		if _, ok := packageDependencyMap[v.PackageName]; !ok {
			packageDependencyMap[v.PackageName] = make([]dependency.Dependency, 0)
		}
		depends := packageDependencyMap[v.PackageName]
		depends = append(depends, v)
		packageDependencyMap[v.PackageName] = depends
	}

	for k, v := range packageDependencyMap {
		depPkgs := map[string]bool{}
		for _, dep := range v {
			for _, depPkg := range dep.DependPkgs {
				depPkgs[depPkg] = true
			}
		}

		for depPkg := range depPkgs {
			_, err := fmt.Fprintf(&b, "%s .> %s\n", k, depPkg)
			if err != nil {
				return err
			}
		}
	}

	b.WriteString("\n")
	b.WriteString("@enduml")

	_, err := w.Write([]byte(b.String()))
	return err
}
