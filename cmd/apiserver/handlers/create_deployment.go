package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	models "github.com/microsoft/commercial-marketplace-offer-deploy/internal/generated"
	"gorm.io/gorm"
)

func CreateDeployment(w http.ResponseWriter, r *http.Request, d data.Database) {
	var command models.CreateDeployment
	err := json.NewDecoder(r.Body).Decode(&command)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validateNameIsUnique(command, w, d.Instance())

	deployment := saveDeployment(command, w, d.Instance())
	utils.WriteJson(w, deployment)
}

func saveDeployment(command models.CreateDeployment, w http.ResponseWriter, db *gorm.DB) *models.Deployment {
	deployment := models.Deployment{
		Name: command.Name,
	}

	tx := db.Create(&deployment)

	if tx.Error != nil {
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return nil
	}
	return &deployment
}

func validateNameIsUnique(command models.CreateDeployment, w http.ResponseWriter, db *gorm.DB) {
	var deployment *models.Deployment
	tx := db.Where("name = ?", command.Name).First(&deployment)

	if tx.Error != nil {
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	alreadyExists := deployment != nil
	if alreadyExists {
		http.Error(w, fmt.Sprintf("Deployment with name [%s] already exists", *deployment.Name), http.StatusBadRequest)
		return
	}
}
