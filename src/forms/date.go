package forms

import "time"

type DateRange struct {
	FromDate int64 `json:"fromDate" form:"fromDate"`
	ToDate   int64 `json:"toDate" form:"toDate"`
}

type DateTimeRange struct {
	FromDate time.Time `json:"fromDate" form:"fromDate"`
	ToDate   time.Time `json:"toDate" form:"toDate"`
}

type RelativeCurrency struct {
	Symbol string `json:"currency" form:"currency"`
}

type TimeInterval struct {
	Interval string `json:"interval" form:"interval"`
}

func (dateRange *DateRange) ToDateTimeRange() *DateTimeRange {
	return &DateTimeRange{
		FromDate: time.Unix(dateRange.FromDate, 0),
		ToDate:   time.Unix(dateRange.ToDate, 0),
	}
}

func (dateTimeRange *DateTimeRange) ToDateRange() *DateRange {
	return &DateRange{
		FromDate: dateTimeRange.FromDate.Unix(),
		ToDate:   dateTimeRange.ToDate.Unix(),
	}
}
