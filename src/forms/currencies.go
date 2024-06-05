package forms

type CurrencyResponse struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	ID     uint   `json:"id"`
}

type CurrenciesResponse struct {
	*SuccessResponse
	Currencies []*CurrencyResponse `json:"currencies"`
}

func (CurrenciesResponse) ListField() string {
	return "Currencies"
}
