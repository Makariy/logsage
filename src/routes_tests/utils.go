package routes_tests

import (
	"main/models"
	"net/url"
	"strconv"
)

func renderDateRangeForTransactions(
	startTransaction *models.Transaction,
	stopTransaction *models.Transaction,
) string {
	params := url.Values{}
	params.Add(
		"fromDate",
		strconv.FormatInt(startTransaction.Date.Unix(), 10),
	)
	params.Add(
		"toDate",
		strconv.FormatInt(stopTransaction.Date.Unix(), 10),
	)
	return params.Encode()
}
