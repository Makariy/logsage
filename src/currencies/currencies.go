package currencies

import "github.com/shopspring/decimal"

// TODO: elaborate mechanism to fetch/set currencies equivalents
func getCurrencyEquivalent(symbol string) decimal.Decimal {
	return decimal.NewFromFloat(1)
}
