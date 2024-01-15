package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/grippenet/helman/pkg/types"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

var (
	ErrFormatNotHandler = errors.New("config file format not handled")
	ErrNoConfigFound    = errors.New("no config file found")
)

func loadFromFile(file string) (*types.Config, error) {
	config := &types.Config{}
	data, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}
	ftype, err := GetConfigType(file)

	if ftype == "yaml" {
		err = yaml.Unmarshal(data, config)
	}
	if ftype == "toml" {
		err = toml.Unmarshal(data, config)
	}
	if err != nil {
		return config, err
	}

	vars, err := resolveVars(config.Vars)
	if err != nil {
		return config, err
	}
	config.Vars = vars

	config.File = file
	return config, nil
}

func resolveVars(vars map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for name, value := range vars {
		if strings.HasPrefix(value, "env:") {
			envName := value[5:]
			envValue := os.Getenv(envName)
			if envValue == "" {
				return nil, fmt.Errorf("environment variable '%s' is not defined or empty for variable '%s'", envName, name)
			}
		} else {
			out[name] = value
		}
	}
	return out, nil
}

func fileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

var defaultPaths = []string{"helman.yaml", "helman.toml"}

func LoadConfig() (*types.Config, error) {
	paths := make([]string, 0)
	file := os.Getenv("HELMAN_CONFIG")
	if file != "" {
		paths = append(paths, file)
	}
	paths = append(paths, defaultPaths...)
	for _, file := range paths {
		ex, err := fileExists(file)
		if err != nil {
			return nil, err
		}
		if ex {
			return loadFromFile(file)
		}
	}
	return nil, ErrNoConfigFound
}

func GetConfigType(file string) (string, error) {
	ext := strings.ToLower(filepath.Ext(file))
	if ext == ".yaml" || ext == ".yml" {
		return "yaml", nil
	}
	if ext == ".toml" {
		return "toml", nil
	}
	return "", ErrFormatNotHandler
}
