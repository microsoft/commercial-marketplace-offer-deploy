package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type UrlHealthCheckOptions struct {
	Timeout           time.Duration
	ReadinessFilePath string
	Url               string
}

type serviceUrlHealthCheck struct {
	options *UrlHealthCheckOptions
}

// Check whether the url is accessible
func (c *serviceUrlHealthCheck) Check(ctx context.Context) HealthCheckResult {
	threshold := time.Now().Add(c.options.Timeout)

	for {
		if time.Now().After(threshold) {
			return HealthCheckResult{
				Description: fmt.Sprintf("Timeout exceeded while waiting for a response from %s", c.options.Url),
				Status:      HealthCheckStatusUnhealthy,
				Error:       errors.New("timeout exceeded while waiting for 200 OK from service url"),
			}
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

func (c *serviceUrlHealthCheck) getResult(ctx context.Context) HealthCheckResult {
	statusCode := 0

	log.Debugf("Checking URL [%s]", c.options.Url)
	response, err := http.Get(c.options.Url)
	if err != nil {
		log.Errorf("Url Health Check error [%s]: %s", c.options.Url, err.Error())
		return HealthCheckResult{
			Description: fmt.Sprintf("Error while checking %s", c.options.Url),
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	statusCode = response.StatusCode

	if statusCode == http.StatusOK {
		c.makeReady()

		return HealthCheckResult{
			Description: fmt.Sprintf("Successfully connected to %s", c.options.Url),
			Status:      HealthCheckStatusHealthy,
		}
	}

	return HealthCheckResult{
		Description: fmt.Sprintf("Received %d from %s", statusCode, c.options.Url),
		Status:      HealthCheckStatusUnhealthy,
		Error:       errors.New("received non-200 response from service url"),
	}
}

func (c *serviceUrlHealthCheck) makeReady() error {
	readiness, err := os.Create(c.options.ReadinessFilePath)

	if err != nil {
		return err
	}
	defer readiness.Close()

	return nil
}

func NewUrlHealthCheck(options UrlHealthCheckOptions) HealthCheck {

	return &serviceUrlHealthCheck{
		options: &options,
	}
}
