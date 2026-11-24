package controller

import (
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/controller/contract"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/controller/gin_adapter"
)

func Default() contract.Controller {
	return gin_adapter.NewController()
}
