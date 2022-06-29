package db

import (
	"errors"
	"fmt"

	"github.com/TestardR/seller-payout/internal/domain"
)

// FindAllSellerWithUnpaidoutItems finds all sellers with unpaid out items.
func (d database) FindSellersWhereItems(where map[string]interface{}) ([]domain.Seller, error) {
	var s []domain.Seller

	db, err := d.preloadSellersRelations(where)
	if err != nil {
		return nil, fmt.Errorf("failed to preload Items: %w", err)
	}

	err = db.FindAll(&s)
	if err == nil {
		return s, nil
	}

	if errors.Is(err, ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return nil, err
}

func (d database) preloadSellersRelations(where map[string]interface{}) (DB, error) {
	tx := d.driver.Preload("Items", where)

	return &database{driver: tx}, tx.Error
}
