package forms

import "time"

type DateRange struct {
	FromDate time.Time `json:"fromDate" form:"fromDate"`
	ToDate   time.Time `json:"toDate" form:"toDate"`
}
