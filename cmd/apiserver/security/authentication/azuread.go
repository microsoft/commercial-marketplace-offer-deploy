package authentication

import (
	"context"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
)

const (
	// JWKS for Azure AD
	// Source: https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration
	AzureAdJwksUri = "https://login.microsoftonline.com/common/discovery/v2.0/keys"

	azureAdTenantIdToken = "{tenantId}"
	AzureAdIssuerFormat  = "https://login.microsoftonline.com/" + azureAdTenantIdToken + "/v2.0"
)

// fetchs the Azure AD key set
func FetchAzureADKeySet(ctx context.Context) (jwk.Set, error) {
	keySet, err := jwk.Fetch(ctx, AzureAdJwksUri)
	return keySet, err
}

func GetAzureAdIssuer(tenantId string) string {
	return strings.Replace(AzureAdIssuerFormat, azureAdTenantIdToken, tenantId, 1)
}
