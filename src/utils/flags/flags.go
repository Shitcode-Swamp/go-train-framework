package flags

import "flag"

type Flags struct {
	Root   *string
	Output *string
	Config *string
}

func NewFlags() *Flags {
	root := flag.String("root", "", "Root path")
	output := flag.String("ctxPath", "./", "Ctx path")
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
