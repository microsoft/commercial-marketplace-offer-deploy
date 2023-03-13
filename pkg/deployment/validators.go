package deployment

type DryRunValidator interface {
	Validate(azureDeployment *AzureDeployment) *DryRunResponse
}

type WhatIfValidatorFunc func(azureDeployment *AzureDeployment) *DryRunResponse

func (f WhatIfValidatorFunc) Validate(azureDeployment *AzureDeployment) *DryRunResponse {
	return f(azureDeployment)
}