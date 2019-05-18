package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/laqiiz/pkgdiagram/dependency"
	"github.com/laqiiz/pkgdiagram/directory"
	"github.com/laqiiz/pkgdiagram/umlgen"

	"github.com/spf13/cobra"
)

var outputPath string

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "file", "f", ".", "output file path")
}

var rootCmd = &cobra.Command{
	Use:   "pgkdgm",
	Short: "pkgdgm is a package diagram generator",
	Long:  "pkgdgm is a tool to analyze package dependencies of Go repository and generate UML text of package diagrams.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0]
		dirName := filepath.Base(targetPath)

		if strings.HasPrefix(targetPath, "https://") {

			prev, _ := filepath.Abs(".")
			defer func() { _ = os.Chdir(prev) }()

			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			tempDir := filepath.Join(userHomeDir, ".pkgdgm")

			// change working directory
			importPath := strings.TrimPrefix(targetPath, "https://")

			downloadPath := filepath.Join(tempDir, importPath)

			// check already imported target repository because go get is very slowly
			if !Exists(downloadPath) {

				parentDir := filepath.Clean(downloadPath + "/../")
				if err := os.MkdirAll(parentDir, os.ModePerm); err != nil && !os.IsExist(err) {
					return err
				}

				// cd
				_ = os.Chdir(parentDir)

				// if not existing then go get
				output, err := exec.Command("git", "clone", targetPath).CombinedOutput()
				if err != nil {
					return err
				}
				log.Println(string(output))
			}

			// Rewrite
			targetPath = filepath.Join(tempDir, importPath)
		}

		dirs, err := directory.NewWithIgnores(".git", ".idea", "docs", "examples").Do(targetPath)
		if err != nil {
			return err
		}

		for _, dir := range dirs {
			log.Printf("%+v\n", dir)
		}

		parser := dependency.New(dirName)

		dependencies := make(dependency.Dependencies, 0)
		for _, v := range dirs {
			parsed := parser.Do(filepath.Join(targetPath, v))
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
			return err
		}

		if err := umlgen.New().Do(file, depends); err != nil {
			return err
		}

		log.Println("finish")

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
