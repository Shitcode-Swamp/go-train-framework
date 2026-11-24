package grpcgen

import (
	"fmt"
	"os"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/grpcgen/codegen"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/grpcgen/registry"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/utils/config"
)

type Builder struct {
	cfg config.GrpcConfig
}

func NewBuilder(cfg config.GrpcConfig) *Builder {
	return &Builder{cfg: cfg}
}

func (b *Builder) Build() error {
	protoRoot, err := b.fetchProtoRoot()
	if err != nil {
		return err
	}
	defer os.RemoveAll(protoRoot)

	if err := codegen.Generate(protoRoot, b.cfg.Output); err != nil {
		return fmt.Errorf("generating grpc code: %w", err)
	}
	return nil
}

func (b *Builder) fetchProtoRoot() (string, error) {
	protoRoot, err := registry.Fetch(b.cfg.SourceRegistry)
	if err != nil {
		return "", fmt.Errorf("fetching proto registry: %w", err)
	}
	return protoRoot, nil
}
