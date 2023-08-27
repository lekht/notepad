package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   `mapstructure:",squash"`
		Postgres `mapstructure:",squash"`
		Logger   `mapstructure:",squash"`
		Speller  `mapstructure:",squash"`
	}

	Server struct {
		Port string `mapstructure:"SERVER_PORT"`
	}

	Postgres struct {
		PgUser     string `mapstructure:"PG_USER"`
		PgPassword string `mapstructure:"PG_PASSWORD"`
		PgHost     string `mapstructure:"PG_HOST"`
		PgPort     string `mapstructure:"PG_PORT"`
		PgDB       string `mapstructure:"PG_DB"`
	}

	Logger struct {
		Level string `mapstructure:"LOGGER_LEVEL"`
	}

	Speller struct {
		Address  string `mapstructure:"SPELLER_ADDRESS"`
		Attempts int    `mapstructure:"SPELLER_ATTEMPTS"`
		Timeout  int64  `mapstructure:"SPELLER_TIMEOUT"`

		Langs   []string `mapstructure:"SPELLER_LANG"`
		Options string   `mapstructure:"SPELLER_OPTION"`
		Format  string   `mapstructure:"SPELLER_FORMAT"`
	}
)

func (c *Config) LoadEnv(path string) error {

	if path == "" {
		files, err := os.ReadDir(".")
		if err != nil {
			return errors.Wrap(err, "failed to find config")
		}

		for _, file := range files {

			filename := file.Name()

			if ext := filepath.Ext(filename); ext != ".env" {
				continue
			}

			if err := c.load("./" + filename); err != nil {
				return errors.Wrap(err, "failed to load config")
			}
			return nil
		}

	}

	if err := c.load(path); err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	return nil
}

func (c *Config) load(path string) error {
	dir, file := filepath.Split(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(file)
	name := filename[0 : len(filename)-len(ext)]

	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(name)
	v.SetConfigType(ext[1:])
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	err = v.Unmarshal(&c)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config to struct")
	}

	return nil
}
