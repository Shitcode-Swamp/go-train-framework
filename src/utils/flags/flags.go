package flags

import (
	"flag"

	"github.com/Shitcode-Swamp/go-train-framework/src/utils/config"
)

type Flags struct {
	Root   *string
	Output *string
	Config *string
}

func NewFlags() *Flags {
	root := flag.String("root", config.DefaultRootPath, "Root path")
	output := flag.String("ctxPath", config.DefaultOutputPath, "Context output path")
	cfg := flag.String("config", "", "Optional config path")
	return (&Flags{
		Root:   root,
		Output: output,
		Config: cfg,
	}).Parse()
}

func (f *Flags) Parse() *Flags {
	flag.Parse()
	return f
}
