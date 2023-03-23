package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/utils"
	data "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
)

func CreateDeployment(w http.ResponseWriter, r *http.Request, d data.Database) {
	var command *generated.CreateDeployment
	err := json.NewDecoder(r.Body).Decode(&command)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deployment := data.FromCreateDeployment(command)

	log.Printf("deployment mapped %v", deployment)

	tx := d.Instance().Create(&deployment)

	if tx.Error != nil {
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, deployment)
}
