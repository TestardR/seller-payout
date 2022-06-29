package handler

import (
	"errors"
	"fmt"

	"github.com/TestardR/seller-payout/internal/domain"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/shopspring/decimal"
)

var errRecoverFromPanic = errors.New("panic defer handler")

const totalPriceLimit = 1_000_000

// CreatePayouts is a background task which goal is to create payouts.
func (h Handler) CreatePayouts() error {
	h.Log.Info("payouts creation started")

	sellers, err := h.DB.FindSellersWhereItems(map[string]interface{}{"paid_out": false})
	if err != nil {
		err = fmt.Errorf("%w: %s", errDB, err)
		h.Log.Error(err)

		return err
	}

	var currencies []domain.Currency
	if err := h.DB.FindAll(&currencies); err != nil {
		err = fmt.Errorf("%w: %s", errDB, err)
		h.Log.Error(err)

		return err
	}

	currenciesMap := make(map[string]domain.Currency)
	for _, c := range currencies {
		currenciesMap[c.Code] = c
	}

	for _, seller := range sellers {
		if len(seller.Items) == 0 {
			continue
		}
		// Concurrent Pipeline organizing payouts creation stages
		if err := h.setupPipeline(seller, currenciesMap); err != nil {
			h.Log.Error(err)

			return err
		}
	}

	h.Log.Info("payouts creation finished")

	return nil
}

// setupPipeline organizes stages for staged processing.
func (h Handler) setupPipeline(seller domain.Seller, currenciesMap map[string]domain.Currency) error {
	// if an error occurs the done channel will gracefully terminate stages 1. and 2.
	done := make(chan struct{})
	defer close(done)

	// Stage 1. creates batch of items
	itemsBatchC := generateItemsBatch(done, seller, currenciesMap)
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
	items      []domain.Item
	totalPrice decimal.Decimal
}

func generateItemsBatch(
	done <-chan struct{},
	seller domain.Seller,
	currencies map[string]domain.Currency) <-chan itemsBatch {
	itemsBatchC := make(chan itemsBatch)

	go func() {
		defer close(itemsBatchC)

		var batch []domain.Item

		totalPrice := decimal.NewFromInt(0)

		for _, item := range seller.Items {
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
	seller domain.Seller,
	currencies map[string]domain.Currency,
	itemsBatchC <-chan itemsBatch) <-chan domain.Payout {
	payoutC := make(chan domain.Payout)
	sellerCurrency := seller.CurrencyCode

	go func() {
		defer close(payoutC)

		for batch := range itemsBatchC {
			p := domain.Payout{
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

func (h Handler) persistPayouts(payoutC <-chan domain.Payout) error {
	runTransaction := func(payout domain.Payout) error {
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
	currencies map[string]domain.Currency,
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
