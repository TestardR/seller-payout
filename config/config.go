package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
)

var errParseEnv = errors.New("failed to parse environment variable")

// Conf represents the application configuration.
type Conf struct {
	// App config
	Port string `required:"true"`
	Env  string `required:"true" validate:"eq=debug|eq=release"`
	CronIntervals
	// Postgres config
	PGUser     string `required:"true" split_words:"true"`
	PGName     string `required:"true" split_words:"true"`
	PGPassword string `required:"true" split_words:"true"`
	PGHost     string `required:"true" split_words:"true"`
}

// CronIntervals represents the intervals from conjobs.
type CronIntervals struct {
	PayoutInterval   int `required:"true" split_words:"true"`
	CurrencyInterval int `required:"true" split_words:"true"`
}

// New returns a new instance of Conf struct.
func New() (Conf, error) {
	var c Conf

	if err := envconfig.Process("", &c); err != nil {
		return Conf{}, fmt.Errorf("%w: %s", errParseEnv, err)
	}

	if err := validator.New().Struct(&c); err != nil {
		return Conf{}, fmt.Errorf("%w: %s", errParseEnv, err)
	}

	return c, nil
}
