package main

import (
	"errors"
	"fmt"

	"github.com/TestardR/seller-payout/pkg/http"
	"github.com/TestardR/seller-payout/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

const appName = "seller-payout"

var (
	errParseEnv      = errors.New("failed to parse environment variable")
	errStoreInstance = errors.New("failed to instanciate storage")
)

type conf struct {
	// App config
	Port string `required:"true"`
	Env  string `required:"true" validate:"eq=debug|eq=release"`
	// Postgres config
	/* PGUser     string `required:"true" split_words:"true"`
	PGName     string `required:"true" split_words:"true"`
	PGPassword string `required:"true" split_words:"true"`
	PGHost     string `required:"true" split_words:"true"` */

	// Redis config
	/* RedisHost string `required:"true" split_words:"true"`
	RedisPort string `required:"true" split_words:"true"` */
}

func config() (conf, error) {
	var c conf

	if err := envconfig.Process("", &c); err != nil {
		return conf{}, fmt.Errorf("%w: %s", errParseEnv, err)
	}

	if err := validator.New().Struct(&c); err != nil {
		return conf{}, fmt.Errorf("%w: %s", errParseEnv, err)
	}

	return c, nil
}

func main() {
	log := logger.NewLogger(appName)

	c, err := config()
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServer(c.Env, log)

	err = server.Run(":" + c.Port)
	if err != nil {
		log.Fatal(err)
	}

}
