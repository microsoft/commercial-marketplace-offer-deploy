package hosting

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
)

type SecurityContext struct {
}

var (
	securityMutex   sync.Mutex
	securityContext *SecurityContext
)

func (c *SecurityContext) GetAzureCredential() azcore.TokenCredential {
	// this will allow us to control how we construct the credential
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	// we MUST have a credential or we're dead in the water anyway
	if err != nil {
		log.Errorf("failed to create credential: %v", err)
	}
	return credential
}

func GetAzureCredentialFunc() func() azcore.TokenCredential {
	context := getDefaultContext()
	return context.GetAzureCredential
}

func GetAzureCredential() azcore.TokenCredential {
	context := getDefaultContext()
	return context.GetAzureCredential()
}

func CheckRoleAssignmentsForScope(appConfig *config.AppConfig, scope string, roleDefinition string, duration time.Duration) (bool, error) {
	threshold := time.Now().Add(duration)

	for {
		if time.Now().After(threshold) {
			return false, errors.New("timeout waiting for role assignment")
		}
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Printf("Failed to obtain a credential: %v", err)
			continue
		}
		ctx := context.Background()
		//todo: use the credential function if necessary
		clientFactory, err := armauthorization.NewClientFactory(appConfig.Azure.SubscriptionId, cred, nil)
		if err != nil {
			log.Printf("Error creating client factory: %v", err)
			continue
		}
		pager := clientFactory.NewRoleAssignmentsClient().NewListForScopePager(scope, nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				log.Printf("Error pulling from pager in NewListForScopePager: %v", err)
				break
			}
			for _, v := range page.Value {
				if v.Properties.PrincipalID == &appConfig.Azure.ClientId && *v.Properties.RoleDefinitionID == roleDefinition {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func getDefaultContext() *SecurityContext {
	securityMutex.Lock()
	defer securityMutex.Unlock()
	if securityContext == nil {
		securityContext = &SecurityContext{}
	}
	return securityContext
}
