package db

import (
	"errors"
	"fmt"

	"github.com/TestardR/seller-payout/internal/domain"
)

// Conditions helper used for Storager queries (see gorm.Where).
type Conditions map[string]interface{}

// FindPayoutsBySellerID finds payouts by seller_id.
func (d database) FindPayoutsBySellerID(id string) ([]domain.Payout, error) {
	where := Conditions{"seller_id": id}

	return d.payouts(where)
}

func (d database) payouts(where Conditions) ([]domain.Payout, error) {
	var p []domain.Payout

	db, err := d.preloadPayoutsRelations()
	if err != nil {
		return nil, fmt.Errorf("failed to preload Items: %w", err)
	}

	err = db.FindAllWhere(&p, where)
	if err == nil {
		return p, nil
	}

	if errors.Is(err, ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return nil, err
}

func (d database) preloadPayoutsRelations() (DB, error) {
	tx := d.driver.Preload("Currency").Preload("Items")

	return &database{driver: tx}, tx.Error
}
