package configer

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-rath/osv-detector/internal/reporter"
	"github.com/g-rath/osv-detector/pkg/database"
	"gopkg.in/yaml.v3"
)

type rawDatabaseConfig struct {
	URL              string `yaml:"url"`
	Name             string `yaml:"name,omitempty"`
	Type             string `yaml:"type,omitempty"`
	WorkingDirectory string `yaml:"working-directory,omitempty"`
}

type rawConfig struct {
	FilePath  string              `yaml:"-"`
	Ignore    []string            `yaml:"ignore"`
	Databases []rawDatabaseConfig `yaml:"extra-databases,omitempty"`
}

type Config struct {
	FilePath  string
	Ignore    []string          `yaml:"ignore"`
	Databases []database.Config `yaml:"extra-databases"`
}

func (rdbc rawDatabaseConfig) inferDBType() string {
	switch {
	case rdbc.Type != "":
		return rdbc.Type
	case strings.HasPrefix(rdbc.URL, "file:/"):
		return "dir"
	case strings.HasSuffix(rdbc.URL, ".zip"):
		return "zip"
		// assume that the url is meant to be for an api,
		// which we will validate later
	default:
		return "api"
	}
}

func (rdbc rawDatabaseConfig) toConfig() (database.Config, error) {
	// the url should always be valid, even for "file" (which should start with "file://")
	_, err := url.ParseRequestURI(rdbc.URL)

	if err != nil {
		return database.Config{}, fmt.Errorf("bad database source url: %w", err)
	}

	finalType := rdbc.inferDBType()

	if finalType != "zip" && finalType != "api" && finalType != "dir" {
		return database.Config{}, fmt.Errorf("%w %s", database.ErrUnsupportedDatabaseType, finalType)
	}

	config := database.Config{
		Name:             rdbc.Name,
		Type:             finalType,
		URL:              rdbc.URL,
		WorkingDirectory: rdbc.WorkingDirectory,
	}

	if config.Name == "" {
		config.Name = config.Identifier()
	}

	return config, nil
}

func newConfig(r *reporter.Reporter, raw rawConfig) (Config, error) {
	config := Config{FilePath: raw.FilePath, Ignore: raw.Ignore, Databases: make(
		[]database.Config,
		0,
		len(raw.Databases),
	)}

	for _, d := range raw.Databases {
		dbc, err := d.toConfig()

		if err != nil {
			r.PrintErrorf("%s contains an invalid database: %v\n", raw.FilePath, err)

			continue
		}

		config.Databases = append(config.Databases, dbc)
	}

	return config, nil
}

// Find attempts to locate & load a Config using the default name (".osv-detector")
func Find(r *reporter.Reporter, pathToDirectory string) (Config, error) {
	var config Config
	var err error

	configName := ".osv-detector"

	config, err = Load(r, pathToDirectory+"/"+configName+".yml")

	if err == nil {
		return config, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return config, err
	}

	config, err = Load(r, pathToDirectory+"/"+configName+".yaml")

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

func load(pathToConfig string) (rawConfig, error) {
	var raw rawConfig

	pathToConfig = filepath.Clean(pathToConfig)

	raw.FilePath = pathToConfig

	configContents, err := os.ReadFile(pathToConfig)

	if err != nil {
		return raw, fmt.Errorf("could not read %s: %w", pathToConfig, err)
	}

	err = yaml.Unmarshal(configContents, &raw)

	if err != nil {
		return raw, fmt.Errorf("could not read %s: %w", pathToConfig, err)
	}

	return raw, nil
}

func Load(r *reporter.Reporter, pathToConfig string) (Config, error) {
	raw, err := load(pathToConfig)

	if err != nil {
		return Config{FilePath: raw.FilePath}, err
	}

	return newConfig(r, raw)
}
