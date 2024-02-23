package forms

type SuccessResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type ListResponse interface {
	ListField() string
}

var (
	Success = &SuccessResponse{
		Status: "success",
	}
)
