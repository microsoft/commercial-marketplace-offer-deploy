package hosting

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
	ctx := context.Background()
	for {
		if time.Now().After(threshold) {
			return false, errors.New("timeout waiting for role assignment")
		}
		cred, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			log.Printf("Failed to obtain a credential: %v", err)
			continue
		}

		accessToken, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{"https://management.azure.com/.default"}})
		if err != nil {
			log.Errorf("failed to get token: %v", err)
			continue
		}
		log.Infof("Token: %s", accessToken.Token)

		if err != nil {
			log.Printf("Failed to obtain a credential: %v", err)
			continue
		}

		objectId, err := getObjectId(&accessToken.Token)
		if err != nil {
			log.Errorf("failed to get object id from token: %v", err)
			continue
		}

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
				if strings.EqualFold(*v.Properties.PrincipalID, objectId) && strings.EqualFold(*v.Properties.RoleDefinitionID, roleDefinition) {
					return true, nil
				} 
			}
		}
	}
}

func getDefaultContext() *SecurityContext {
	securityMutex.Lock()
	defer securityMutex.Unlock()
	if securityContext == nil {
		securityContext = &SecurityContext{}
	}
	return securityContext
}

func getObjectId(rawToken *string) (string, error) {
	token, err := jwt.Parse(*rawToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keySet, err := FetchAzureADKeySet(context.Background())
		if err != nil {
			return nil, fmt.Errorf("could not fetch keyset")
		}
		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}
		return publickey, nil
	})
	if err != nil {
		return "", err
	}
	return token.Claims.(jwt.MapClaims)["oid"].(string), nil
}
