package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App           App
	Observability Observability
}

type App struct {
	Name string
	Port string
}

type Observability struct {
	OtelEndpoint string
}

func Init() (*Config, error) {
	configPath := "./order_service/config/config.yaml"

	cfg := &Config{}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}
