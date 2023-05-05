package config

import (
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Domain       string `json:"domain" yaml:"domain"`
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
}

func NewEmptyConfig() *Config {
	return &Config{}
}

// Loads the config from a file.
func NewConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return &Config{}, nil
	}

	exists, err := fileExists(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine if the config file %s exists", configPath)
	}

	if !exists {
		return nil, errors.Errorf("Config file %s does not exist", configPath)
	}

	v := viper.New()
	v.SetConfigFile("yaml")
	v.AddConfigPath(".")
	v.SetConfigFile(configPath)
	v.SetEnvPrefix("DS_LOAD_AUTH0")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err = v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open config file '%s'", configPath)
	}

	v.AutomaticEnv()

	cfg := new(Config)
	err = v.UnmarshalExact(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config file '%s'", configPath)
	}

	return cfg, nil
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "failed to stat file '%s'", path)
	}
}
