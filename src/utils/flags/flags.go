package flags

import "flag"

type Flags struct {
	Root       *string
	Output     *string
}

func NewFlags() Flags {
	root := flag.String("root", "", "root path")
	output := flag.String("ctxPath", "./", "ctx path")
	return Flags{
		Root: root,
		Output: output,
	}
}

func (f *Flags) Parse() {
	flag.Parse()
}
