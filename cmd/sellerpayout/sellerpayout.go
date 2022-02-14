package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/TestardR/seller-payout/internal/handler"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

const (
	appName      = "seller-payout"
	migrationDir = "migrations"
)

var errParseEnv = errors.New("failed to parse environment variable")

type conf struct {
	// App config
	Port             string `required:"true"`
	Env              string `required:"true" validate:"eq=debug|eq=release"`
	PayoutInterval   int    `required:"true" split_words:"true"`
	CurrencyInterval int    `required:"true" split_words:"true"`
	// Postgres config
	PGUser     string `required:"true" split_words:"true"`
	PGName     string `required:"true" split_words:"true"`
	PGPassword string `required:"true" split_words:"true"`
	PGHost     string `required:"true" split_words:"true"`
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
	log := logger.New(appName)

	c, err := config()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.New(db.Config{
		User:     c.PGUser,
		Name:     c.PGName,
		Password: c.PGPassword,
		Host:     c.PGHost,
	})
	if err != nil {
		log.Fatal("failed to create postgres: %w", err)
	}

	if err = db.RunMigrations(migrationDir); err != nil {
		log.Fatal("failed to run postgres migration: %w", err)
	}

	h := handler.Handler{
		Log: log,
		DB:  db,
		EX:  currency.New(),
	}

	payoutTicker := time.NewTicker(time.Duration(c.PayoutInterval) * time.Hour)
	currencyTicker := time.NewTicker(time.Duration(c.CurrencyInterval) * time.Hour)

	go func() {
		for {
			select {
			case <-payoutTicker.C:
				if err := h.CreatePayouts(); err != nil {
					log.Fatal("failed to create payouts: %w", err)
				}
			case <-currencyTicker.C:
				if err := h.UpdateCurrencies(); err != nil {
					log.Fatal("failed to update currencies: %w", err)
				}
			}
		}
	}()

	server := h.NewServer(c.Env)

	err = server.Run(":" + c.Port)
	if err != nil {
		log.Fatal(err)
	}
}
