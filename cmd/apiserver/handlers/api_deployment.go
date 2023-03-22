package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	models "github.com/microsoft/commercial-marketplace-offer-deploy/internal/generated"
)

func CreateDeployment(w http.ResponseWriter, r *http.Request, d data.Database) {
	var command models.CreateDeployment
	err := json.NewDecoder(r.Body).Decode(&command)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := d.Instance()
	var deployment models.Deployment
	deployment = db.Where("name = ?", command.Name).First(&deployment)

	if deployment != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(command)
}

func GetDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func InvokeOperation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ListDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
