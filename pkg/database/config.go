package database

import (
	"errors"
	"fmt"
)

type Config struct {
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	URL              string `yaml:"url"`
	WorkingDirectory string `yaml:"working-directory"`
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

var ErrUnsupportedDatabaseType = errors.New("unsupported database source type")

// Load initializes a new OSV database based on the given Config
func Load(config Config, offline bool, batchSize int, pkgNames []string) (DB, error) {
	switch config.Type {
	case "zip":
		return NewZippedDB(config, offline, pkgNames)
	case "api":
		return NewAPIDB(config, offline, batchSize)
	case "dir":
		return NewDirDB(config, offline, pkgNames)
	}

	return nil, fmt.Errorf("%w %s", ErrUnsupportedDatabaseType, config.Type)
}
