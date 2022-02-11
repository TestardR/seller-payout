package currency

import (
	"errors"
	"fmt"

	"github.com/asvvvad/exchange"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=currency.go -destination=$MOCK_FOLDER/currency.go -package=mock

const (
	// USDCode is the USD currency.
	USDCode = "USD"
	// GBPCode is the GBP currency.
	GBPCode = "GBP"
	// EURCode is the EUR currency.
	EURCode = "EUR"
)

var supportedCurrency = map[string]struct{}{
	USDCode: {},
	GBPCode: {},
	EURCode: {},
}

var (
	errInvalidCurrency  = errors.New("currency is not supported")
	errExchangeAPI      = errors.New("an error occurred with the external API")
	errConvertToDecimal = errors.New("failed converting to decimal")
)

// Exchanger is the currency exchange interface.
type Exchanger interface {
	// GetConversionRate takes in a currency and returns its exchange rate against 1 unit of in the base currency.
	GetConversionRate(currency string) (decimal.Decimal, error)
}

type exchanger struct {
	api *exchange.Exchange
}

// New returns a new client instance with base currency set by default to USD.
func New() Exchanger {
	return exchanger{
		api: exchange.New(USDCode),
	}
}

func (e exchanger) GetConversionRate(currency string) (decimal.Decimal, error) {
	if _, ok := supportedCurrency[currency]; !ok {
		return decimal.Decimal{}, fmt.Errorf("%w (format: %s, accepted: %v)", errInvalidCurrency, currency, supportedCurrency)
	}

	rate, err := e.api.ConvertTo(currency, 1)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("%w:%s", errExchangeAPI, err)
	}

	dec, err := decimal.NewFromString(rate.String())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("%w:%s", errConvertToDecimal, err)
	}

	return dec, nil
}
