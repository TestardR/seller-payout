package http

import (
	"fmt"

	"github.com/TestardR/seller-payout/internal/domain"
)

// UpdateCurrencies is a background task,
// it calls an external API to update currencies exchange rate.
func (h Handler) UpdateCurrencies() error {
	h.Log.Info("currencies update started")

	var currencies []domain.Currency

	if err := h.DB.FindAll(&currencies); err != nil {
		err = fmt.Errorf("%w: %s", errDB, err)
		h.Log.Error(err)

		return err
	}

	for i, c := range currencies {
		rate, err := h.EX.GetConversionRate(c.Code)
		if err != nil {
			h.Log.Error(err)

			return err
		}

		currencies[i].USDExchRate = rate
	}

	if err := h.DB.Update(currencies); err != nil {
		err = fmt.Errorf("%w: %s", errDB, err)
		h.Log.Error(err)

		return err
	}

	h.Log.Info("currencies update started")

	return nil
}
