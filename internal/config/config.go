package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OutputDir string `mapstructure:"output_dir"`
	Port      string `mapstructure:"port"`
}

func Load() (*Config, error) {
	viper.SetConfigName("api-replay")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	return &cfg, nil
}
