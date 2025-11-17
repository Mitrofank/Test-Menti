package models

import "fmt"

type CurrencyCode string

const (
	CurrencyRUB CurrencyCode = "RUB"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyEUR CurrencyCode = "EUR"
)

var supportedCurrencies = map[CurrencyCode]struct{}{
	CurrencyRUB: {},
	CurrencyUSD: {},
	CurrencyEUR: {},
}

func IsSupported(code string) bool {
	_, ok := supportedCurrencies[CurrencyCode(code)]
	return ok
}

func (c CurrencyCode) Validate() error {
	if !IsSupported(string(c)) {
		return fmt.Errorf("unsupported currency code: %s", c)
	}
	return nil
}
