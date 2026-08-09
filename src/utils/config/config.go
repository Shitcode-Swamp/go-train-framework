package config

import (
	"path/filepath"

	"github.com/spf13/viper"
	"fmt"
)

type grpcRegistry struct {
	SourceRegistry string `yaml:"sourceRegistry" mapstructure:"sourceRegistry"`
	IsGitHub       bool   `yaml:"isGithub" mapstructure:"isGithub"`
}

type grpcConfig struct {
	SourceRegistry grpcRegistry `yaml:"sourceRegistry" mapstructure:"sourceRegistry"`
	Output         string       `yaml:"output" mapstructure:"output"`
}

type Config struct {
	Root   string `yaml:"root" mapstructure:"root"`
	Output string `yaml:"output" mapstructure:"output"`

	GrpcConfig grpcConfig `yaml:"grpcConfig" mapstructure:"grpcConfig"`
}

func NewConfig(root, output string) *Config {
	return &Config{
		Root:   root,
		Output: output,
		GrpcConfig: grpcConfig{},
	}
}

func ConfigFromFile(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(fileNameFromPath(path))
	v.SetConfigType(fileExtFromPath(path))
	v.AddConfigPath(filePathWithoutName(path))
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}


func (c *Config) OverrideRoot(root string) *Config {
	c.Root = root
	return c
}

func (c *Config) OverrideOutput(output string) *Config {
	c.Output = output
	return c
}

func fileNameFromPath(path string) string {
	return filepath.Base(path)
}

func fileExtFromPath(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		return ext[1:]
	}
	return ""
}

func filePathWithoutName(path string) string {
	return filepath.Dir(path)
}

func LogConfig(c *Config) {
	fmt.Printf("Config:\n")
	fmt.Printf("  Root: %s\n", c.Root)
	fmt.Printf("  Output: %s\n", c.Output)
	fmt.Printf("  GrpcConfig:\n")
	fmt.Printf("    SourceRegistry:\n")
	fmt.Printf("      SourceRegistry: %s\n", c.GrpcConfig.SourceRegistry.SourceRegistry)
	fmt.Printf("      IsGitHub: %t\n", c.GrpcConfig.SourceRegistry.IsGitHub)
	fmt.Printf("    Output: %s\n", c.GrpcConfig.Output)
}
