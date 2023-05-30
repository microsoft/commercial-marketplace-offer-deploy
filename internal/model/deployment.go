package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"gorm.io/gorm"
)

const deploymentNamePartSeparator = "."

type Deployment struct {
	gorm.Model
	Name     string         `json:"name"`
	Template map[string]any `json:"template" gorm:"json"`
	Stages   []Stage        `json:"stages" gorm:"json"`

	// azure properties
	SubscriptionId string `json:"subscriptionId"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
}

// Gets the azure deployment name suitable for azure deployment
// format - modm.<deploymentId>.<deploymentName>
func (d *Deployment) GetAzureDeploymentName() string {
	prefix := "modm" + deploymentNamePartSeparator + strconv.FormatUint(uint64(d.ID), 10) + deploymentNamePartSeparator
	suffix := d.getSanitizedName()
	maxLength := 64
	lengthCheck := len(prefix + suffix)

	// reduce size of the name if the total length is greater than 64
	if lengthCheck > maxLength {
		suffix = suffix[:maxLength-len(suffix)]
	}
	return prefix + suffix
}

// Parses the azure deployment name and returns OUR deployment id
func (d *Deployment) ParseAzureDeploymentName(azureDeploymentName string) (*int, error) {
	isManagedDeployment := strings.HasPrefix(azureDeploymentName, deployment.LookupPrefix)

	if isManagedDeployment {
		parts := strings.Split(azureDeploymentName, deploymentNamePartSeparator)
		idString := parts[1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			return nil, err
		}
		return &id, nil
	}
	return nil, fmt.Errorf("[%s] is not a managed deployment", azureDeploymentName)
}

func (d *Deployment) getSanitizedName() string {
	r := regexp.MustCompile("[^a-zA-Z0-9 -]")
	name := r.ReplaceAllString(d.Name, "")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.TrimSuffix(name, "-")

	return name
}
