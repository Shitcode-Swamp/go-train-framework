package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/grpcgen"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/config"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/flags"
)

func main() {
	f := flags.NewGrpcFlags()

	cfg, err := loadGrpcConfig(*f.Config)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	if err := grpcgen.NewBuilder(cfg).Build(); err != nil {
		fmt.Println("error: ", err)
	}
}

func loadGrpcConfig(path string) (config.GrpcConfig, error) {
	cfg, err := config.ConfigFromFile(path)
	if err != nil {
		return config.GrpcConfig{}, err
	}
	return cfg.GrpcConfig, nil
}
