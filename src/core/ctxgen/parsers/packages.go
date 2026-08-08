package parsers

import (
	"golang.org/x/tools/go/packages"
)

type Packages struct {
	pattern  string
	packages map[string]string
}

func NewPackages(pattern string) *Packages {
	packages, err := getPackages(pattern)
	if err != nil {
		panic("Error parsing packages")
	}
	return &Packages{
		pattern:  pattern,
		packages: packages,
	}
}

func (p *Packages) GetPackage(filePath string) (string, bool) {
	if p.packages == nil {
		return "", false
	}
	val, ok := p.packages[filePath]
	return val, ok
}

func getPackages(pattern string) (map[string]string, error) {
	config := packages.Config{Mode: packages.NeedName | packages.NeedFiles}
	packages, err := packages.Load(&config, pattern)
	if err != nil {
		return nil, err
	}
	packagesMap := make(map[string]string, len(packages))
	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			packagesMap[file] = pkg.PkgPath
		}
	}
	return packagesMap, nil
}
