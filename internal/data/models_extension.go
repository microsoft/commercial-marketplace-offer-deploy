package data

// The purpose of this file is to provide a place to put extension methods to the data models
// so we keep models clean

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

// Gets the azure deployment name suitable for azure deployment
// format - modm.<deploymentId>-<deploymentName>
func (d *Deployment) GetAzureDeploymentName() string {
	return "modm." + strconv.FormatUint(uint64(d.ID), 10) + "-" + d.Name
}

// Parses the azure deployment name and returns OUR deployment id
func (d *Deployment) ParseAzureDeploymentName(azureDeploymentName string) (*int, error) {
	parts := strings.Split(azureDeploymentName, "-")
	isManagedDeployment := strings.HasPrefix(parts[0], deployment.LookupPrefix)

	if isManagedDeployment {
		idString := strings.TrimPrefix(parts[0], deployment.LookupPrefix)
		id, err := strconv.Atoi(idString)
		if err != nil {
			return nil, err
		}
		return &id, nil
	}
	return nil, fmt.Errorf("[%s] is not a managed deployment", azureDeploymentName)
}

//region InvokedOperation

func (io *InvokedOperation) IsRetriable() bool {
	return io.Retries > io.Attempts
}

//endregion InvokedOperation
