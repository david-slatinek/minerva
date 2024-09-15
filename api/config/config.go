package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ConnectionString string `mapstructure:"connection-string"`
	Mode             string `mapstructure:"mode"`
}

func NewConfig(filename string) (*Config, error) {
	cfg := &Config{}
	err := cfg.loadConfig(filename)

	if err != nil {
		return nil, err
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (receiver *Config) loadConfig(filename string) error {
	viper.SetConfigName(filename)
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&receiver)
}

func (receiver *Config) validate() error {
	if receiver.Mode == "" {
		log.Println("mode not set, setting it to release")
		receiver.Mode = "release"
	}

	if receiver.ConnectionString == "" {
		return errors.New("connection-string not set")
	}

	return nil
}
