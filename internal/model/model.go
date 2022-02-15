package model

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

// Seller is an individual owning items.
type Seller struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	CurrencyCode string `json:"currency_code"`

	Items []Item
}

// Item is a sold product.
type Item struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	ReferenceName string          `json:"reference_name"`
	PriceAmount   decimal.Decimal `json:"price_amount"`
	CurrencyCode  string          `json:"currency_code"`
	PaidOut       bool            `json:"-"`

	// https://gorm.io/docs/belongs_to.html#Belongs-To
	SellerID uuid.UUID `gorm:"type:uuid" json:"seller_id"`
	Seller   Seller    `gorm:"foreignKey:seller_id" json:"seller"`
}

// Currency is the money in which transactions are made.
type Currency struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Code        string          `json:"code"`
	USDExchRate decimal.Decimal `json:"usd_exch_rate"`
}
