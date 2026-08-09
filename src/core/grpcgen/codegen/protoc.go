package codegen

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

var ErrNoProtoFiles = errors.New("no .proto files found in source registry")

func Generate(protoRoot, output string) error {
	protoFiles, err := findProtoFiles(protoRoot)
	if err != nil {
		return err
	}
	if len(protoFiles) == 0 {
		return ErrNoProtoFiles
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	return runProtoc(protoRoot, output, protoFiles)
}

func findProtoFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(path) == ".proto" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk %s: %w", root, err)
	}
	return files, nil
}

func runProtoc(protoRoot, output string, protoFiles []string) error {
	args := append(protocArgs(protoRoot, output), protoFiles...)

	cmd := exec.Command("protoc", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("protoc failed: %w: %s", err, stderr.String())
	}
	return nil
}

func protocArgs(protoRoot, output string) []string {
	return []string{
		"--proto_path=" + protoRoot,
		"--go_out=" + output,
		"--go_opt=paths=source_relative",
		"--go-grpc_out=" + output,
		"--go-grpc_opt=paths=source_relative",
	}
}
