package umlgen

import (
	"fmt"
	"github.com/laqiiz/pkgdiagram/dependency"
	"io"
	"strings"
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

	// package
	for _, v := range dependencies {
		_, err := fmt.Fprintf(&b, "package %s{}\n", v.PackageName)
		if err != nil {
			return err
		}
	}
	b.WriteString("\n")

	// link
	for _, v := range dependencies {

		for _, dependPkg := range v.DependPkgs {
			_, err := fmt.Fprintf(&b, "%s .> %s\n", v.PackageName, dependPkg)
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
