package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/appctx"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/storegen"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/utils/config"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/utils/flags"
)

func main() {
	f := flags.NewFlags()
	cfg, err := config.NewConfig(*f.Root, *f.Output, *f.Config)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	config.LogConfig(cfg)
	if err := storegen.NewStoreFieldsBuilder(cfg.Root).Build(); err != nil {
		fmt.Println("error: ", err)
		return
	}
	appBuilder := appctx.NewAppCtxBuilder(cfg.Root, cfg.Output)
	if err := appBuilder.Build(); err != nil {
		fmt.Println("error: ", err)
	}
}
