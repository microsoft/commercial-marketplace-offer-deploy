package deployment

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type DryRunValidator interface {
	Validate(input DryRunValidationInput) *DryRunResponse
}

type WhatIfValidatorFunc func(input DryRunValidationInput) *DryRunResponse

func (f WhatIfValidatorFunc) Validate(input DryRunValidationInput) *DryRunResponse {
	return f(input)
}

type DryRunValidationInput struct {
	ctx context.Context
	cred azcore.TokenCredential
	azureDeployment *AzureDeployment
}