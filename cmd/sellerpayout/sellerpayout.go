package main

import (
	"time"

	"github.com/TestardR/seller-payout/config"
	"github.com/TestardR/seller-payout/internal/handler/http"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/logger"
)

const (
	appName      = "seller-payout"
	migrationDir = "migrations"
)

func main() {
	log := logger.New(appName)

	c, err := config.New()
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

	h := http.Handler{
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
