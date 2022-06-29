package db

import (
	"errors"
	"fmt"

	"github.com/TestardR/seller-payout/internal/domain"
)

// FindUnpaidOutItemsBySellerID finds unpaid out itmes by seller_id.
func (d database) FindUnpaidOutItemsBySellerID(id string) ([]domain.Item, error) {
	where := Conditions{"seller_id": id, "paid_out": false}

	return d.items(where)
}

// FindUnpaidOutItems finds unpaid out itmes.
func (d database) FindUnpaidOutItems() ([]domain.Item, error) {
	where := Conditions{"paid_out": false}

	return d.items(where)
}

func (d database) items(where Conditions) ([]domain.Item, error) {
	var items []domain.Item

	db, err := d.preloadItemsRelations()
	if err != nil {
		return nil, fmt.Errorf("failed to preload Currencies: %w", err)
	}

	err = db.FindAllWhere(&items, where)
	if err == nil {
		return items, nil
	}

	if errors.Is(err, ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return nil, err
}

func (d database) preloadItemsRelations() (DB, error) {
	tx := d.driver.Preload("Seller")

	return &database{driver: tx}, tx.Error
}
