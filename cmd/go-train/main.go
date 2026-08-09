package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/ctxgen/appctx"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/config"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/flags"
)

func main() {
	f := flags.NewFlags()
	cfg, err := config.NewConfig(*f.Root, *f.Output, *f.Config)
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
