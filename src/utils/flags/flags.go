package flags

import "flag"

const (
	DefaultRootPath   = ""
	DefaultOutputPath = "./"
)

type Flags struct {
	Root   *string
	Output *string
	Config *string
}

func NewFlags() *Flags {
	root := flag.String("root", DefaultRootPath, "Root path")
	output := flag.String("output", DefaultOutputPath, "Output path")
	config := flag.String("config", "", "Optional config path")
	return &Flags{
		Root:   root,
		Output: output,
		Config: config,
	}
}

func (f *Flags) Parse() {
	flag.Parse()
}
