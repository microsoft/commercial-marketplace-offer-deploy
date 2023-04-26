package hosting

import (
	"context"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
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
