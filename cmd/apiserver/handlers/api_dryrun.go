package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

func CreateDryRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var azureDeployment deployment.AzureDeployment
	err := json.NewDecoder(r.Body).Decode(&azureDeployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// write response in bondy indicating a failed marshall
		return
	}
	res := deployment.DryRun(&azureDeployment)
	respondJSON(w, http.StatusOK, res)
}