package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error: Unable to write error to http output stream: %v", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if status != http.StatusOK {
		// Avoids log message about setting header superfluously
		w.WriteHeader(status)
	}
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Printf("Error: Unable to write response to http output stream: %v", err)
	}
}