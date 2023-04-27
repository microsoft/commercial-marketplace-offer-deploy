package diagnostics

import "context"

const (
	HealthCheckStatusHealthy   = 2
	HealthCheckStatusDegraded  = 1
	HealthCheckStatusUnhealthy = 0
)

type HealthCheckResult struct {
	Description string
	Status      int
	Error       error
}

// A health check, which can be used to check the status of a component in the application, such as a backend service, database or some internal state.
type HealthCheck interface {
	Check(ctx context.Context) HealthCheckResult
}

type HealthCheckService interface {
	AddHealthCheck(healthCheck HealthCheck)
	// Runs the provided health checks and returns the aggregated status.
	CheckHealth(ctx context.Context) []HealthCheckResult
}

type healthCheckService struct {
	checks []HealthCheck
}

// AddHealthCheck implements HealthCheckService
func (service *healthCheckService) AddHealthCheck(check HealthCheck) {
	service.checks = append(service.checks, check)
}

func (service *healthCheckService) CheckHealth(ctx context.Context) []HealthCheckResult {
	results := make([]HealthCheckResult, len(service.checks))
	for i, check := range service.checks {
		results[i] = check.Check(ctx)
	}
	return results
}

func NewHealthCheckService() HealthCheckService {
	return &healthCheckService{
		checks: []HealthCheck{},
	}
}
