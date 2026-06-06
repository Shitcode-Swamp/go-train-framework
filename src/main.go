package main

import (
	"fmt"

	files_traverse "github.com/Shitcode-Swamp/go-train/src/core/files_traverse"
	// "github.com/Shitcode-Swamp/go-train/src/core/injection"
	// "github.com/Shitcode-Swamp/go-train/src/core/someshi"
	"github.com/Shitcode-Swamp/go-train/src/utils/flags"
)

func main() {
	flags := flags.NewFlags()
	flags.Parse()
	fmt.Println("root =", *flags.Root)
	entries, err := files_traverse.Parse(*flags.Root)
	if err != nil {
		fmt.Println("error: ", err)
	}
	for token, entry := range entries.Injections {
		fmt.Printf("Token: %s,\n", token)
		entry.Expose()
	}
}