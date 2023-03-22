package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	datamodel "github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
)

func CreateDeployment(w http.ResponseWriter, r *http.Request, d data.Database) {
	var command internal.CreateDeployment
	err := json.NewDecoder(r.Body).Decode(&command)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deployment := datamodel.FromCreateDeployment(&command)

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
