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
	appBuilder := appctx.NewAppCtxBuilder(cfg.Root, cfg.Output)
	if err := appBuilder.Build(); err != nil {
		fmt.Println("error: ", err)
	}
}

func getConfig(flags *flags.Flags) (*config.Config, error) {
	if *flags.Config == "" {
		return config.NewConfig(*flags.Root, *flags.Output), nil
	}
	cfg, err := config.ConfigFromFile(*flags.Config)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}
	return cfg.
		OverrideRoot(*flags.Root).
		OverrideOutput(*flags.Output), nil
}
