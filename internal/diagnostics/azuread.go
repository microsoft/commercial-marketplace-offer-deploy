package diagnostics

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	log "github.com/sirupsen/logrus"
)

const (
	// JWKS for Azure AD
	// Source: https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration
	AzureAdJwksUri = "https://login.microsoftonline.com/common/discovery/v2.0/keys"

	azureAdTenantIdToken  = "{tenantId}"
	AzureAdV2IssuerFormat = "https://login.microsoftonline.com/" + azureAdTenantIdToken + "/v2.0/"
	AzureAdV1IssuerFormat = "https://sts.windows.net/" + azureAdTenantIdToken + "/"
)

// fetchs the Azure AD key set
func FetchAzureADKeySet(ctx context.Context) (jwk.Set, error) {
	keySet, err := jwk.Fetch(ctx, AzureAdJwksUri)
	return keySet, err
}

func GetAzureAdIssuers(tenantId string) []string {
	issuers := []string{
		strings.Replace(AzureAdV2IssuerFormat, azureAdTenantIdToken, tenantId, 1),
		strings.Replace(AzureAdV1IssuerFormat, azureAdTenantIdToken, tenantId, 1),
	}
	return issuers
}

func GetAzureCredentialObjectId(ctx context.Context) (string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", err
	}

	accessToken, err := cred.GetToken(ctx, policy.TokenRequestOptions{})
	if err != nil {
		log.Errorf("failed to get JWT token to extract ObjectID from the Azure Credential: %v", err)
	}

	rawToken := accessToken.Token

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keySet, err := FetchAzureADKeySet(ctx)
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
