package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train/src/core/appctx"
	"github.com/Shitcode-Swamp/go-train/src/utils/flags"
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
