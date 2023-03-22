package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
)

type InvokeOperationDeploymentHandler func(internal.InvokeDeploymentOperation, data.Database) (interface{}, error)

// func CreateDeployment(w http.ResponseWriter, r *http.Request, d data.Database) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)
// }

func GetDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func InvokeOperation(w http.ResponseWriter, r *http.Request, d data.Database) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var operation internal.InvokeDeploymentOperation
	err := json.NewDecoder(r.Body).Decode(&operation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	operationHandler := CreateOperationHandler(operation)
	if operationHandler == nil {
		http.Error(w, "There was op OperationHandler registered for the invoked operation", http.StatusBadRequest)
	}
	
	res, err := operationHandler(operation, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	respondJSON(w, http.StatusOK, res)
}

func CreateOperationHandler(operation internal.InvokeDeploymentOperation) InvokeOperationDeploymentHandler {
	switch operation.Name {
	case operation.Name:
		return CreateDryRun
	default:
		return nil
	}
}

func ListDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
