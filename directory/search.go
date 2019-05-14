package directory

import (
	"os"
	"path/filepath"
	"strings"
)

type Searcher struct {
	ignores []string
}

func NewWithIgnores(ignores ...string) *Searcher {
	return &Searcher{
		ignores: ignores,
	}
}

func (s *Searcher) Ignore(str string) bool {
	for _, v := range s.ignores {
		if strings.HasPrefix(str, v) {
			return true
		}
	}
	return false
}

func (s *Searcher) Do(basePath string) ([]string, error) {
	dirs := make([]string, 0)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}
			if s.Ignore(relPath) {
				return nil
			}
			dirs = append(dirs, relPath)
		}
		return nil
	})
	return dirs, err
}
