package testutils

import (
	"encoding/json"
	"net/http/httptest"
)

func NewFromResponseRecorder[T any](response *httptest.ResponseRecorder) (T, error) {
	var v T
	err := json.Unmarshal(response.Body.Bytes(), &v)
	return v, err
}
