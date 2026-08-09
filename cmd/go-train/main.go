package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/ctxgen/appctx"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/config"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/flags"
)

func main() {
	flags := flags.NewFlags()
	flags.Parse()
	cfg, err := getConfig(flags)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	config.LogConfig(cfg)
	appBuilder := appctx.NewAppCtxBuilder(cfg.Root, cfg.Output)
	if err := appBuilder.Build(); err != nil {
		fmt.Println("error: ", err)
	}
}

func getConfig(f *flags.Flags) (*config.Config, error) {
	if *f.Config == "" {
		return config.NewConfig(*f.Root, *f.Output), nil
	}
	cfg, err := config.ConfigFromFile(*f.Config)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}
	if *f.Root != flags.DefaultRootPath {
		cfg.OverrideRoot(*f.Root)
	}
	if *f.Output != flags.DefaultOutputPath {
		cfg.OverrideOutput(*f.Output)
	}
	return cfg, nil
}
