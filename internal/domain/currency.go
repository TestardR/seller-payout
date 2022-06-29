package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// Currency is the money in which transactions are made.
type Currency struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Code        string          `json:"code"`
	USDExchRate decimal.Decimal `json:"usd_exch_rate"`
}
