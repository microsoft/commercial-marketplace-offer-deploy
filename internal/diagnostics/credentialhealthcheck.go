package diagnostics

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/gommon/log"
)

type AzureCredentialHealthCheckOptions struct {
	Timeout time.Duration
}

type azureCredentialHealthCheck struct {
	options *AzureCredentialHealthCheckOptions
}

// Check implements diagnostics.HealthCheck
func (c *azureCredentialHealthCheck) Check(ctx context.Context) HealthCheckResult {
	threshold := time.Now().Add(c.options.Timeout)

	for {
		if time.Now().After(threshold) {
			return HealthCheckResult{
				Description: "Timeout exceeded while waiting for role assignments to be created",
				Status:      HealthCheckStatusUnhealthy,
				Error:       errors.New("timeout exceeded while waiting for role assignments to be created"),
			}
		}
		result := c.getResult(ctx)

		if result.Status != HealthCheckStatusHealthy || result.Error != nil {
			log.Warnf("Health Check attempt failed: %v", result)
			continue
		}

		return result
	}
}

func (c *azureCredentialHealthCheck) getResult(ctx context.Context) HealthCheckResult {
	_, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return HealthCheckResult{
			Description: "Failed to get Azure credential",
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	return HealthCheckResult{
		Description: "Azure Credential Health Check",
		Status:      HealthCheckStatusHealthy,
	}
}

func NewAzureCredentialHealthCheck(options AzureCredentialHealthCheckOptions) HealthCheck {

	return &azureCredentialHealthCheck{
		options: &options,
	}
}
