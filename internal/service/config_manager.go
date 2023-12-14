package service

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"dough-calculator/internal/config"
	"dough-calculator/internal/domain"
)

type configManager struct {
	config config.Config
}

func (configManager *configManager) ParseConfig() error {
	configPath := os.Getenv("CONFIG_PATH")

	configFile, err := os.Open(configPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open config file %s", configPath)
	}
	viperService := viper.NewWithOptions(viper.KeyDelimiter("_"))
	viperService.SetConfigType("yaml")
	viperService.AutomaticEnv()

	err = viperService.ReadConfig(configFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read config file '%s'", configPath)
	}

	err = viperService.Unmarshal(&configManager.config)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config file '%s'", configPath)
	}

	if configManager.config == (config.Config{}) {
		return errors.Errorf("failed to parse config file '%s'", configPath)
	}

	return nil
}

func (configManager *configManager) GetConfig() config.Config {
	return configManager.config
}

func NewConfigManager() domain.ConfigManager {
	return &configManager{}
}
