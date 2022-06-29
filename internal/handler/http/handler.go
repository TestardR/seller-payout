package http

import (
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/logger"
)

type handler struct {
	Log logger.Logger
	DB  db.DB
	EX  currency.Exchanger
}
