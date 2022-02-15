package handler

import (
	"errors"
	"testing"

	"github.com/TestardR/seller-payout/internal/model"
	"github.com/TestardR/seller-payout/pkg/mock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type handleCaseCreatePayouts struct {
	h   Handler
	err error
}

func TestHandler_CreatePayouts(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)
	t.Cleanup(func() { mc.Finish() })

	tests := map[string]handleCaseCreatePayouts{
		"fail-db-find-unpaid-items":     payoutsCreateCaseFailDBFindUnpaidOutItems(mc),
		"fail-db-find-currencies":       payoutsCreateCaseFailDBFindCurrencies(mc),
		"fail-db-begin-tx":              payoutsCreateCaseFailDBBeginTX(mc),
		"fail-db-insert-tx":             payoutsCreateCaseFailDBInsertTX(mc),
		"fail-db-update-tx":             payoutsCreateCaseFailDBUpdateTX(mc),
		"fail-db-commit-tx":             payoutsCreateCaseFailDBCommitTX(mc),
		"split-payouts-above-max-price": payoutsCreateCaseSplitPayoutsAboveMaxPrice(mc),
		"success":                       payoutsCreateCaseOK(mc),
	}

	for tn, tc := range tests {
		tn, tc := tn, tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			err := tc.h.CreatePayouts()

			if (err != nil) != (tc.err != nil) {
				assert.Equal(t, err, tc.err)
			}
		})
	}
}

func payoutsCreateCaseFailDBCommitTX(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	merr := errors.New("mock")

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItems(), nil)
	mdb.EXPECT().FindAll(gomock.Any())
	mdb.EXPECT().Begin().Return(mdb, nil)
	mdb.EXPECT().Insert(gomock.Any())
	mdb.EXPECT().Update(validItems(true))
	mdb.EXPECT().Commit().Return(merr)
	ml.EXPECT().Error(gomock.Any())
	ml.EXPECT().Error(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: merr,
	}
}

func payoutsCreateCaseFailDBUpdateTX(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	merr := errors.New("mock")

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItems(), nil)
	mdb.EXPECT().FindAll(gomock.Any())
	mdb.EXPECT().Begin().Return(mdb, nil)
	mdb.EXPECT().Insert(gomock.Any())
	mdb.EXPECT().Update(validItems(true)).Return(merr)
	ml.EXPECT().Error(gomock.Any())
	ml.EXPECT().Error(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: merr,
	}
}

func payoutsCreateCaseFailDBInsertTX(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	merr := errors.New("mock")

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItems(), nil)
	mdb.EXPECT().FindAll(gomock.Any())
	mdb.EXPECT().Begin().Return(mdb, nil)
	mdb.EXPECT().Insert(gomock.Any()).Return(merr)
	ml.EXPECT().Error(gomock.Any())
	ml.EXPECT().Error(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: merr,
	}
}

func payoutsCreateCaseFailDBBeginTX(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	merr := errors.New("mock")

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItems(), nil)
	mdb.EXPECT().FindAll(gomock.Any())
	mdb.EXPECT().Begin().Return(nil, merr)
	ml.EXPECT().Error(gomock.Any())
	ml.EXPECT().Error(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: merr,
	}
}

func payoutsCreateCaseFailDBFindCurrencies(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	merr := errors.New("mock")

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItems(), nil)
	mdb.EXPECT().FindAll(gomock.Any()).Return(merr)
	ml.EXPECT().Error(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: merr,
	}
}

func payoutsCreateCaseFailDBFindUnpaidOutItems(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	merr := errors.New("mock")

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return([]model.Seller{}, merr)
	ml.EXPECT().Error(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: merr,
	}
}

func payoutsCreateCaseSplitPayoutsAboveMaxPrice(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItemsAboveMaxPrice(), nil)
	mdb.EXPECT().FindAll(gomock.Any())
	mdb.EXPECT().Begin().Return(mdb, nil)
	mdb.EXPECT().Insert(gomock.Any())
	mdb.EXPECT().Update(gomock.Any())
	mdb.EXPECT().Commit()

	mdb.EXPECT().Begin().Return(mdb, nil)
	mdb.EXPECT().Insert(gomock.Any())
	mdb.EXPECT().Update(gomock.Any())
	mdb.EXPECT().Commit()
	ml.EXPECT().Info(gomock.Any())
	ml.EXPECT().Info(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: nil,
	}
}

func payoutsCreateCaseOK(mc *gomock.Controller) handleCaseCreatePayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	ml.EXPECT().Info(gomock.Any())
	mdb.EXPECT().FindSellersWhereItems(map[string]interface{}{"paid_out": false}).Return(validSellersWithUnpaidOutItems(), nil)
	mdb.EXPECT().FindAll(gomock.Any())
	mdb.EXPECT().Begin().Return(mdb, nil)
	mdb.EXPECT().Insert(gomock.Any())
	mdb.EXPECT().Update(validItems(true))
	mdb.EXPECT().Commit()
	ml.EXPECT().Info(gomock.Any())
	ml.EXPECT().Info(gomock.Any())

	return handleCaseCreatePayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		err: nil,
	}
}

func Test_convertToSellerCurrency(t *testing.T) {
	currencies := make(map[string]model.Currency)

	currencies["USD"] = model.Currency{
		USDExchRate: decimal.NewFromInt(1),
	}
	currencies["EUR"] = model.Currency{
		USDExchRate: decimal.NewFromInt(1),
	}

	got := convertToSellerCurrency("USD", "EUR", currencies, decimal.NewFromInt(1))
	assert.True(t, got.Equal(decimal.NewFromInt(1)))

	got = convertToSellerCurrency("EUR", "USD", currencies, decimal.NewFromInt(1))
	assert.True(t, got.Equal(decimal.NewFromInt(1)))
}

func validItems(paidout bool) []model.Item {
	mID := uuid.FromStringOrNil("test-id")
	mItem := model.Item{
		ReferenceName: "test-ref-name",
		SellerID:      mID,
		Seller:        model.Seller{ID: mID, CurrencyCode: "USD"},
		PaidOut:       paidout,
		PriceAmount:   decimal.NewFromInt(1000000),
		CurrencyCode:  "USD",
	}

	return []model.Item{mItem}
}

func validItemsAboveMaxPrice(paidout bool) []model.Item {
	mID := uuid.FromStringOrNil("test-id")
	mItem := model.Item{
		ReferenceName: "test-ref-name",
		SellerID:      mID,
		Seller:        model.Seller{ID: mID, CurrencyCode: "USD"},
		PaidOut:       paidout,
		PriceAmount:   decimal.NewFromInt(1000000),
		CurrencyCode:  "USD",
	}

	return []model.Item{mItem, mItem}
}

func validSellersWithUnpaidOutItems() []model.Seller {
	mSeller := model.Seller{
		CurrencyCode: "USD",
		Items:        validItems(false),
	}

	return []model.Seller{mSeller}
}

func validSellersWithUnpaidOutItemsAboveMaxPrice() []model.Seller {
	return []model.Seller{
		{
			CurrencyCode: "USD",
			Items:        validItemsAboveMaxPrice(false)},
	}
}
