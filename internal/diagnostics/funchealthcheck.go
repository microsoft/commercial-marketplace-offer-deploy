package diagnostics

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type FuncHealthCheckOptions struct {
	Timeout time.Duration
	Ready   func() bool
}

type funcHealthCheck struct {
	options *FuncHealthCheckOptions
}

// Check implements diagnostics.HealthCheck
func (c *funcHealthCheck) Check(ctx context.Context) HealthCheckResult {
	threshold := time.Now().Add(c.options.Timeout)

	for {
		if time.Now().After(threshold) {
			result := HealthCheckResult{
				Description: "Timeout exceeded while waiting for ready result",
				Status:      HealthCheckStatusUnhealthy,
				Error:       errors.New("timeout exceeded while waiting for ready result"),
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

func (c *funcHealthCheck) getResult(ctx context.Context) HealthCheckResult {
	result := c.options.Ready()

	if !result {
		return HealthCheckResult{
			Description: "Ready func returned false",
			Status:      HealthCheckStatusUnhealthy,
			Error:       errors.New("failed readiness check"),
		}
	}
	return HealthCheckResult{
		Description: "Ready func returned true",
		Status:      HealthCheckStatusHealthy,
	}
}

func NewFuncHealthCheck(options FuncHealthCheckOptions) HealthCheck {
	return &funcHealthCheck{
		options: &options,
	}
}
