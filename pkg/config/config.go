package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port   string `env:"PORT,required"`
	DbConn string `env:"DB_CONN,required"`
}

var C Config

func Load() error {
	err := env.Parse(&C)
	if err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	return &C
}
