package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	log "github.com/sirupsen/logrus"
)

type AzureRoleAssignmentsHealthCheckOptions struct {
	SubscriptionId    string
	RoleAssignmentIds []string
	Timeout           time.Duration
}

type azureRoleAssignmentsHealthCheck struct {
	options *AzureRoleAssignmentsHealthCheckOptions
}

// Check implements diagnostics.HealthCheck
func (c *azureRoleAssignmentsHealthCheck) Check(ctx context.Context) HealthCheckResult {
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
			time.Sleep(5 * time.Second)
			continue
		}

		return result
	}
}

func (c *azureRoleAssignmentsHealthCheck) getResult(ctx context.Context) HealthCheckResult {
	objectId, err := GetAzureCredentialObjectId(ctx)
	if err != nil {
		return HealthCheckResult{
			Description: "Failed to get Azure credential object ID",
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}
	roleAssignments, err := c.getRoleAssignments(objectId, ctx)

	if err != nil {
		return HealthCheckResult{
			Description: "Failed to get Azure role assignments",
			Status:      HealthCheckStatusUnhealthy,
			Error:       err,
		}
	}

	if len(roleAssignments) == 0 {
		return HealthCheckResult{
			Description: fmt.Sprintf("No role assignments found for Azure credential object ID %s", objectId),
			Status:      HealthCheckStatusUnhealthy,
			Error:       nil,
		}
	}

	roleAssignmentsToCheckCount := len(c.options.RoleAssignmentIds)
	counter := 0
	for _, roleAssignmentId := range c.options.RoleAssignmentIds {
		if c.isRoleAssigned(roleAssignmentId, roleAssignments) {
			counter++
		}
	}

	if counter == roleAssignmentsToCheckCount {
		return HealthCheckResult{
			Description: fmt.Sprintf("All role assignments found for Azure credential object ID %s", objectId),
			Status:      HealthCheckStatusHealthy,
		}
	}

	return HealthCheckResult{
		Description: fmt.Sprintf("Not all role assignments found for Azure credential object ID %s", objectId),
		Status:      HealthCheckStatusDegraded,
	}
}

func NewRoleAssignmentsHealthCheck(options AzureRoleAssignmentsHealthCheckOptions) HealthCheck {
	return &azureRoleAssignmentsHealthCheck{
		options: &options,
	}
}

func (c *azureRoleAssignmentsHealthCheck) isRoleAssigned(roleAssignmentId string, roleAssignments []*armauthorization.RoleAssignment) bool {
	for _, roleAssignment := range roleAssignments {
		// the role assignment ID is the same as the role assignment name that gets set in the bicep template
		if *roleAssignment.Name == roleAssignmentId {
			return true
		}
	}
	return false
}

func (c *azureRoleAssignmentsHealthCheck) getRoleAssignments(objectId string, ctx context.Context) ([]*armauthorization.RoleAssignment, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armauthorization.NewClientFactory(c.options.SubscriptionId, cred, nil)

	if err != nil {
		return nil, err
	}

	roleAssignmentsClient := clientFactory.NewRoleAssignmentsClient()
	getAllAssignmentsByObjectId := to.Ptr(fmt.Sprintf("principalId eq '%s'", objectId))

	pager := roleAssignmentsClient.NewListForScopePager(&armauthorization.RoleAssignmentsClientListForSubscriptionOptions{
		Filter: getAllAssignmentsByObjectId,
	})

	roleAssignments := []*armauthorization.RoleAssignment{}

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if nextResult.RoleAssignmentListResult.Value != nil {
			roleAssignments = append(roleAssignments, nextResult.RoleAssignmentListResult.Value...)
		}
	}
	return roleAssignments, nil
}
