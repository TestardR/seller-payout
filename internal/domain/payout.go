package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// Payout is an invoice assigned to a seller with a total price in a currency
// for a list of items.
type Payout struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	PriceTotal decimal.Decimal `json:"price_total"`

	// https://gorm.io/docs/belongs_to.html#Belongs-To
	SellerID   uuid.UUID `gorm:"type:uuid" json:"seller_id"`
	Seller     Seller    `gorm:"foreignKey:seller_id" json:"seller"`
	CurrencyID uuid.UUID `gorm:"type:uuid" json:"currency_id"`
	Currency   Currency  `gorm:"foreignKey:currency_id" json:"currency"`

	Items []Item `gorm:"many2many:payout_items;"`
}
