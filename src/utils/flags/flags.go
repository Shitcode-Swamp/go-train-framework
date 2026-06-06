package flags

import "flag"

type Flags struct {
	Root       *string
	AppCtxPath *string
}

func NewFlags() Flags {
	root := flag.String("root", "", "root path")
	appCtxPath := flag.String("ctxPath", "", "ctx path")
	return Flags{
		Root: root,
		AppCtxPath: appCtxPath,
	}
}

func (f *Flags) Parse() {
	flag.Parse()
}
