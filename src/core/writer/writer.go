package writer

import (
	"os"
	"path/filepath"
)

const outputFileName = "app.go"

type FileWriter struct {
	outputDir string
}

func NewFileWriter(outputDir string) *FileWriter {
	return &FileWriter{outputDir: outputDir}
}

func (fw *FileWriter) Write(content string) error {
	if err := os.MkdirAll(fw.outputDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(fw.outputDir, outputFileName), []byte(content), 0o644)
}
