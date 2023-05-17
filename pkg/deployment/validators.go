package deployment

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type DryRunValidator interface {
	Validate(input DryRunValidationInput) (*sdk.DryRunResponse, error)
}

type WhatIfValidatorFunc func(input DryRunValidationInput) (*sdk.DryRunResponse, error)

func (f WhatIfValidatorFunc) Validate(input DryRunValidationInput) (*sdk.DryRunResponse, error) {
	return f(input)
}

type DryRunValidationInput struct {
	ctx             context.Context
	cred            azcore.TokenCredential
	azureDeployment *AzureDeployment
}
