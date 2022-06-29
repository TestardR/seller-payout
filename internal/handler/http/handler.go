package http

import (
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/logger"
)

// Handler is the main structure to inject dependencies.
type Handler struct {
	Log logger.Logger
	DB  db.DB
	EX  currency.Exchanger
}
