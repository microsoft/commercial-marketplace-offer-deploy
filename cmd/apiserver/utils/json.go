package utils

import (
	"encoding/json"
	"net/http"
)

type JsonResponse interface {
	Result(d any)
}

type jsonResponse struct {
	writer http.ResponseWriter
}

func WriteJson(w http.ResponseWriter, d any) {
	response := NewJsonResponse(w)
	response.Result(d)
}

func NewJsonResponse(w http.ResponseWriter) JsonResponse {
	response := &jsonResponse{writer: w}
	return response
}

func (r *jsonResponse) Result(d any) {
	r.writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	r.writer.WriteHeader(http.StatusOK)
	json.NewEncoder(r.writer).Encode(d)
}
