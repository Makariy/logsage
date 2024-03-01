package routes_tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func PerformTestRequest(
	router *gin.Engine,
	method string,
	path string,
	data []byte,
	headers *map[string]string,
) *httptest.ResponseRecorder {
	var req *http.Request
	if data != nil {
		req, _ = http.NewRequest(method, path, strings.NewReader(string(data)))
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
		}
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func GetResponseJSONBody(resp *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var responseBody map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func AssertResponseSuccess(code int, resp *httptest.ResponseRecorder, suite *suite.Suite) {
	suite.Equal(code, resp.Code, "Got an unexpected response code", resp.Body)

	responseBody, err := GetResponseJSONBody(resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal("success", responseBody["status"])
}

func UnmarshalResponse[T any](resp *httptest.ResponseRecorder) (*T, error) {
	var form T
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &form)
	if err != nil {
		return nil, err
	}
	return &form, err
}
