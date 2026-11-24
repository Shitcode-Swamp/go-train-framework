package flags

import (
	"flag"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/utils/config"
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

type GrpcFlags struct {
	Config *string
}

func NewGrpcFlags() *GrpcFlags {
	cfg := flag.String("config", "./go-train.yaml", "Path to the go-train config file")
	return (&GrpcFlags{Config: cfg}).Parse()
}

func (f *GrpcFlags) Parse() *GrpcFlags {
	flag.Parse()
	return f
}
