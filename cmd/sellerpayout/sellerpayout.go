package main

import (
	"github.com/TestardR/seller-payout/config"
	"github.com/TestardR/seller-payout/internal/handler/cron"
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

	cron.Run(log, db, currency.New(), c.CronIntervals)

	server := http.NewServer(c.Env, log, db)

	err = server.Run(":" + c.Port)
	if err != nil {
		log.Fatal(err)
	}
}
