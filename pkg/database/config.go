package database

import (
	"fmt"
)

type Config struct {
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	URL              string `yaml:"url"`
	WorkingDirectory string `yaml:"working-directory"`
}

func NewConfig(name, explicitType, url, workingDirectory string) Config {
	return Config{
		Name:             name,
		Type:             explicitType,
		URL:              url,
		WorkingDirectory: workingDirectory,
	}
}

// Identifier returns a unique string that can be used to check if a loaded
// database has been configured with this Config
func (dbc Config) Identifier() string {
	id := dbc.Type + "#" + dbc.URL

	if dbc.WorkingDirectory != "" {
		id += "#" + dbc.WorkingDirectory
	}

	return id
}

// Load initializes a new OSV database based on the given Config
func Load(config Config, offline bool, batchSize int) (DB, error) {
	switch config.Type {
	case "zip":
		return NewZippedDB(config, offline)
	case "api":
		return NewAPIDB(config, offline, batchSize)
	case "dir":
		return NewDirDB(config, offline)
	}

	return nil, fmt.Errorf("oh noes")
}
