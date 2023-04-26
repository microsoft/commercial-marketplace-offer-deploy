package hosting

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
		log.Fatalf("failed to create credential: %v", err)
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

func getDefaultContext() *SecurityContext {
	securityMutex.Lock()
	defer securityMutex.Unlock()
	if securityContext == nil {
		securityContext = &SecurityContext{}
	}
	return securityContext
}
