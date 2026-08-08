package lexical

import "fmt"

type Import struct {
	Alias string
	Path  string
}

func NewImport(alias, path string) Import {
	return Import{Alias: alias, Path: path}
}

func (i Import) String() string {
	if i.Alias == "" {
		return fmt.Sprintf("%q", i.Path)
	}
	return fmt.Sprintf("%s %q", i.Alias, i.Path)
}
