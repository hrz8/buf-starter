package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/hrz8/altalune"

	"gopkg.in/yaml.v2"
)

var (
	cachedConfig *AppConfig
	once         sync.Once
)

func Load(filename string) (altalune.Config, error) {
	var loadErr error

	once.Do(func() {
		file, err := os.Open(filename)
		if err != nil {
			loadErr = fmt.Errorf("failed to open config file: %w", err)
			return
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)

		cachedConfig = &AppConfig{}
		if err := decoder.Decode(cachedConfig); err != nil {
			loadErr = fmt.Errorf("failed to decode config file: %w", err)
			return
		}

		cachedConfig.setDefaults()

		if err := cachedConfig.Validate(); err != nil {
			loadErr = fmt.Errorf("configuration validation failed: %w", err)
			return
		}
	})

	return cachedConfig, loadErr
}
