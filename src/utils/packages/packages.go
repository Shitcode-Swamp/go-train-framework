package packages

import (
	"fmt"
	"hash/fnv"
)

type Package struct {
	Name string
	Path string
}

func (p Package) Alias() string {
	h := fnv.New32a()
	h.Write([]byte(p.Path))
	return fmt.Sprintf("%s%06x", p.Name, h.Sum32()&0xFFFFFF)
}
