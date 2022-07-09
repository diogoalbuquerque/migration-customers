package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		Log        `yaml:"logger"`
		MYSQL      `yaml:"mysql"`
		DB2        `yaml:"db2"`
		AWS        `yaml:"aws"`
		DATASOURCE `yaml:"datasource"`
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	MYSQL struct {
		OpenConnMax int `env-required:"true" env:"MYSQL_OPEN_CONN_MAX"`
		IdleConnMax int `env-required:"true" env:"MYSQL_IDLE_CONN_MAX"`
		LifeConnMax int `env-required:"true" env:"MYSQL_LIFE_CONN_MAX"`
	}

	DB2 struct {
		OpenConnMax int `env-required:"true" env:"DB2_OPEN_CONN_MAX"`
		IdleConnMax int `env-required:"true" env:"DB2_IDLE_CONN_MAX"`
		LifeConnMax int `env-required:"true" env:"DB2_LIFE_CONN_MAX"`
	}

	AWS struct {
		RegionName string `env-required:"true"  env:"AWS_REGION_NAME"`
		SecretName string `env-required:"true" env:"SECRETS_MANAGER"`
	}

	DATASOURCE struct {
		Limit int `env-required:"true" env:"DATASOURCE_LIMIT"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
