package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/ctxgen/appctx"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/flags"
)

func main() {
	flags := flags.NewFlags()
	flags.Parse()
	fmt.Println("root =", *flags.Root)
	appBuilder := appctx.NewAppCtxBuilder(*flags.Root, *flags.Output)
	if err := appBuilder.Build(); err != nil {
		fmt.Println("error: ", err)
	}
}
