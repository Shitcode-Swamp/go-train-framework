package registry

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/Shitcode-Swamp/go-train-framework/src/utils/config"
)

var ErrMissingSourceRegistry = errors.New("grpcConfig.sourceRegistry.sourceRegistry must be set")

func Fetch(cfg config.GrpcRegistry) (string, error) {
	if cfg.SourceRegistry == "" {
		return "", ErrMissingSourceRegistry
	}

	dir, err := newTempCloneDir()
	if err != nil {
		return "", err
	}

	url := cloneURL(cfg)
	if err := cloneInto(dir, url); err != nil {
		os.RemoveAll(dir)
		return "", err
	}
	return dir, nil
}

func newTempCloneDir() (string, error) {
	dir, err := os.MkdirTemp("", "go-train-grpc-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	return dir, nil
}

func cloneInto(dir, url string) error {
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{URL: url, Depth: 1}); err != nil {
		return fmt.Errorf("failed to clone %s: %w", url, err)
	}
	return nil
}

func cloneURL(cfg config.GrpcRegistry) string {
	src := cfg.SourceRegistry
	if isGitHubShorthand(cfg) {
		return githubURL(src)
	}
	return src
}

func isGitHubShorthand(cfg config.GrpcRegistry) bool {
	return cfg.IsGitHub && !strings.Contains(cfg.SourceRegistry, "://")
}

func githubURL(ownerRepo string) string {
	return "https://github.com/" + strings.TrimSuffix(ownerRepo, ".git") + ".git"
}
