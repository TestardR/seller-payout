package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

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
