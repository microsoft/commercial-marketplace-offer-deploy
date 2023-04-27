package diagnostics

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	log "github.com/sirupsen/logrus"
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
			result := HealthCheckResult{
				Description: "Timeout exceeded while waiting for azure credential check",
				Status:      HealthCheckStatusUnhealthy,
				Error:       errors.New("timeout exceeded while waiting for azure credential check"),
			}
			log.Warnf("Health Check timed out: %v", result)
			return result
		}

		result := c.getResult(ctx)

		if result.Status != HealthCheckStatusHealthy || result.Error != nil {
			log.Warnf("Health Check attempt failed: %v", result)
			time.Sleep(5 * time.Second)
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
