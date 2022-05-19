package configer

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	FilePath string
	Ignore   []string `yaml:"ignore"`
}

// Find attempts to locate & load a Config using the default name (".osv-detector")
func Find(pathToDirectory string) (Config, error) {
	var config Config
	var err error

	configName := ".osv-detector"

	config, err = Load(pathToDirectory + "/" + configName + ".yml")

	if err == nil {
		return config, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return config, err
	}

	config, err = Load(pathToDirectory + "/" + configName + ".yaml")

	if err == nil {
		return config, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return config, err
	}

	// if we couldn't find a config at all,
	// we want to return an empty Config
	// that doesn't have FilePath set
	return Config{}, nil
}

func Load(pathToConfig string) (Config, error) {
	var config Config

	config.FilePath = pathToConfig

	configContents, err := os.ReadFile(pathToConfig)

	if err != nil {
		return config, fmt.Errorf("could not read %s: %w", pathToConfig, err)
	}

	err = yaml.Unmarshal(configContents, &config)

	if err != nil {
		return config, fmt.Errorf("could not read %s: %w", pathToConfig, err)
	}

	return config, nil
}
