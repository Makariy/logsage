package forms

type CurrencyResponse struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

type CurrenciesResponse struct {
	*SuccessResponse
	Currencies []*CurrencyResponse `json:"currencies"`
}

func (CurrenciesResponse) ListField() string {
	return "Currencies"
}
