package configer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func UpdateWithIgnores(pathToConfig string, ignores []string) error {
	raw, err := load(pathToConfig)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	raw.Ignore = ignores

	out, err := yaml.Marshal(raw)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	f, err := os.OpenFile(pathToConfig, os.O_TRUNC|os.O_WRONLY, os.ModePerm)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = f.Write(out)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
