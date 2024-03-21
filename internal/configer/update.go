package configer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func UpdateWithIgnores(pathToConfig string, ignores []string) error {
	raw, err := load(pathToConfig)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	raw.Ignore = ignores

	f, err := os.OpenFile(pathToConfig, os.O_TRUNC|os.O_WRONLY, os.ModePerm)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	err = encoder.Encode(raw)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
