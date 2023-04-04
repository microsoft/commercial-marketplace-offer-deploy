package data

// The purpose of this file is to provide a place to put extension methods to the data models
// so we keep models clean

import (
	"fmt"
	"strconv"
	"strings"
)

const DeploymentPrefix = "modm."

// Gets the azure deployment name suitable for azure deployment
// format - modm.<deploymentId>-<deploymentName>
func (d *Deployment) GetAzureDeploymentName() string {
	return "modm." + strconv.FormatUint(uint64(d.ID), 10) + "-" + d.Name
}

// Parses the azure deployment name and returns the deployment id
func (d *Deployment) ParseAzureDeploymentName(azureDeploymentName string) (int, error) {
	parts := strings.Split(azureDeploymentName, "-")
	isManagedDeployment := strings.HasPrefix(parts[0], DeploymentPrefix)

	if isManagedDeployment {
		id := strings.TrimPrefix(parts[0], DeploymentPrefix)
		return strconv.Atoi(id)
	}
	return -1, fmt.Errorf("[%s] is not a managed deployment", azureDeploymentName)
}
