package lexical

type Package struct {
	Name string
}

func NewPackage(name string) Package {
	return Package{Name: name}
}

func (p Package) String() string {
	return "package " + p.Name + "\n"
}