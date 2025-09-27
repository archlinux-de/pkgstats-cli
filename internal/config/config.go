package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Blocklist struct {
		Packages []string `mapstructure:"packages"`
		Mirrors  []string `mapstructure:"mirrors"`
	} `mapstructure:"blocklist"`
}

func parse(v *viper.Viper) (*Config, error) {
	var config Config

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return &config, nil
		} else {
			return nil, fmt.Errorf("failed to read config file %s: %w", v.ConfigFileUsed(), err)
		}
	}

	if err := v.UnmarshalExact(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", v.ConfigFileUsed(), err)
	}

	return &config, nil
}

func Load(configFile string) (*Config, error) {
	v := viper.New()

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.AddConfigPath("/etc")
		v.SetConfigName("pkgstats")
		v.SetConfigType("yaml")
	}

	return parse(v)
}
