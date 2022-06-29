package cron

import (
	"errors"
	"time"

	"github.com/TestardR/seller-payout/config"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/logger"
)

var (
	errCreatePayouts    = errors.New("failed to create payouts")
	errUpdateCurrencies = errors.New("failed to update currencies")
)

type handler struct {
	Log logger.Logger
	DB  db.DB
	EX  currency.Exchanger
}

// Run initializes cron jobs.
func Run(log logger.Logger, db db.DB, ex currency.Exchanger, c config.CronIntervals) {
	h := handler{
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
					log.Fatal("%w: %s", errCreatePayouts, err)
				}
			case <-currencyTicker.C:
				if err := h.UpdateCurrencies(); err != nil {
					log.Fatal("%w: %s", errUpdateCurrencies, err)
				}
			}
		}
	}()
}
