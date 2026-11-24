package storegen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type StoreFieldsBuilder struct {
	root string
}

func NewStoreFieldsBuilder(root string) *StoreFieldsBuilder {
	return &StoreFieldsBuilder{root: root}
}

func (b *StoreFieldsBuilder) Build() error {
	root, err := filepath.Abs(b.root)
	if err != nil {
		return err
	}
	return filepath.WalkDir(root, func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !isSourceFile(filePath) {
			return nil
		}
		return b.generate(filePath)
	})
}

func (b *StoreFieldsBuilder) generate(sourcePath string) error {
	file, err := parseEntityFile(sourcePath)
	if err != nil {
		return err
	}
	if !file.hasEntities() {
		return nil
	}
	formatted, err := format.Source([]byte(file.render()))
	if err != nil {
		return fmt.Errorf("format store fields for %s: %w", sourcePath, err)
	}
	return writeGenerated(file.outputPath(), formatted)
}

func isSourceFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".go") && !strings.HasSuffix(filePath, "_test.go")
}

func writeGenerated(outputPath string, content []byte) error {
	existing, err := os.ReadFile(outputPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err == nil && !bytes.HasPrefix(existing, []byte(generatedHeader)) {
		return fmt.Errorf("refusing to overwrite hand-written file %s", outputPath)
	}
	return os.WriteFile(outputPath, content, 0o644)
}
