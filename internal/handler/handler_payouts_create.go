package handler

import (
	"errors"
	"fmt"

	"github.com/TestardR/seller-payout/internal/model"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

var errRecoverFromPanic = errors.New("panic defer handler")

const totalPriceLimit = 1_000_000

// CreatePayouts is a background task which goal is to create payouts.
func (h Handler) CreatePayouts() error {
	h.Log.Info("payouts creation started")

	items, err := h.DB.FindUnpaidOutItems()
	if err != nil {
		err = fmt.Errorf("%w: %s", errDB, err)
		h.Log.Error(err)

		return err
	}

	sellersMap := make(map[uuid.UUID]model.Seller)
	for _, item := range items {
		sellersMap[item.SellerID] = item.Seller
	}

	itemsMap := itemsMapFromSeller(items)

	var currencies []model.Currency
	if err := h.DB.FindAll(&currencies); err != nil {
		err = fmt.Errorf("%w: %s", errDB, err)
		h.Log.Error(err)

		return err
	}

	currenciesMap := make(map[string]model.Currency)
	for _, c := range currencies {
		currenciesMap[c.Code] = c
	}

	for sellerID, seller := range sellersMap {
		// Concurrent Pipeline organizing payouts creation stages
		if err := h.setupPipeline(sellerID, seller, itemsMap, currenciesMap); err != nil {
			h.Log.Error(err)

			return err
		}
	}

	h.Log.Info("payouts creation finished")

	return nil
}

// setupPipeline organizes stages for staged processing.
func (h Handler) setupPipeline(
	sellerID uuid.UUID,
	seller model.Seller,
	itemsMap map[uuid.UUID][]model.Item,
	currenciesMap map[string]model.Currency) error {
	// if an error occurs the done channel will terminate stages 1. and 2.
	done := make(chan struct{})
	defer close(done)

	// Stage 1. creates a batch of items
	itemsBatchC := generateItemsBatch(done, seller, itemsMap[sellerID], currenciesMap)
	// Stage 2. creates payouts
	payoutC := generatePayouts(done, seller, currenciesMap, itemsBatchC)
	// Stage 3. persists payouts
	if err := h.persistPayouts(payoutC); err != nil {
		h.Log.Error(err)

		return err
	}

	return nil
}

type itemsBatch struct {
	items      []model.Item
	totalPrice decimal.Decimal
}

func generateItemsBatch(
	done <-chan struct{},
	seller model.Seller,
	items []model.Item,
	currencies map[string]model.Currency) <-chan itemsBatch {
	itemsBatchC := make(chan itemsBatch)

	go func() {
		defer close(itemsBatchC)

		var batch []model.Item

		totalPrice := decimal.NewFromInt(0)

		for _, item := range items {
			price := convertToSellerCurrency(
				seller.CurrencyCode,
				item.CurrencyCode,
				currencies,
				item.PriceAmount)

			if totalPrice.Add(price).GreaterThan(decimal.NewFromInt(totalPriceLimit)) {
				itemsBatchC <- itemsBatch{
					batch,
					totalPrice,
				}

				batch = nil
				totalPrice = decimal.NewFromInt(0)
			}

			totalPrice = totalPrice.Add(price)

			batch = append(batch, item)
		}

		ib := itemsBatch{
			batch,
			totalPrice,
		}

		select {
		case itemsBatchC <- ib:
		case <-done:
		}
	}()

	return itemsBatchC
}

func generatePayouts(
	done <-chan struct{},
	seller model.Seller,
	currencies map[string]model.Currency,
	itemsBatchC <-chan itemsBatch) <-chan model.Payout {
	payoutC := make(chan model.Payout)
	sellerCurrency := seller.CurrencyCode

	go func() {
		defer close(payoutC)

		for batch := range itemsBatchC {
			p := model.Payout{
				PriceTotal: batch.totalPrice.Round(numberOfDecimals),
				Items:      batch.items,
				SellerID:   seller.ID,
				Seller:     seller,
				CurrencyID: currencies[sellerCurrency].ID,
				Currency:   currencies[sellerCurrency],
			}

			select {
			case payoutC <- p:
			case <-done:
			}
		}
	}()

	return payoutC
}

func (h Handler) persistPayouts(payoutC <-chan model.Payout) error {
	runTransaction := func(payout model.Payout) error {
		tx, err := h.DB.Begin()
		if err != nil {
			err = fmt.Errorf("failed to create DB transaction: %w", err)

			return err
		}

		defer func() {
			if r := recover(); r != nil {
				err = tx.Rollback()

				err = errRecoverFromPanic
			}
		}()

		if err := tx.Insert(&payout); err != nil {
			return fmt.Errorf("%w: %s", errDB, err)
		}

		for i := 0; i < len(payout.Items); i++ {
			payout.Items[i].PaidOut = true
		}

		if err := tx.Update(payout.Items); err != nil {
			return fmt.Errorf("%w: %s", errDB, err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction in DB: %w", err)
		}

		return nil
	}

	for payout := range payoutC {
		err := runTransaction(payout)
		if err != nil {
			return err
		}
	}

	h.Log.Info("payout successfully persisted")

	return nil
}

func convertToSellerCurrency(
	sellerCode, itemCode string,
	currencies map[string]model.Currency,
	price decimal.Decimal) decimal.Decimal {
	if itemCode == sellerCode {
		return price
	}

	if sellerCode == currency.USDCode {
		price = price.Div(currencies[itemCode].USDExchRate)
	} else {
		price = price.Div(currencies[itemCode].USDExchRate).Mul(currencies[sellerCode].USDExchRate)
	}

	return price
}

func itemsMapFromSeller(items []model.Item) map[uuid.UUID][]model.Item {
	output := make(map[uuid.UUID][]model.Item)

	for _, item := range items {
		output[item.SellerID] = append(output[item.SellerID], item)
	}

	return output
}
