package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

// Seller is an individual owning items.
type Seller struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	CurrencyCode string `json:"currency_code"`

	Items []Item
}
