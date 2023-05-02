package deployment

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type DryRunValidator interface {
	Validate(input DryRunValidationInput) (*DryRunResponse, error)
}

type WhatIfValidatorFunc func(input DryRunValidationInput) (*DryRunResponse, error)

func (f WhatIfValidatorFunc) Validate(input DryRunValidationInput) (*DryRunResponse, error) {
	return f(input)
}

type DryRunValidationInput struct {
	ctx context.Context
	cred azcore.TokenCredential
	azureDeployment *AzureDeployment
}