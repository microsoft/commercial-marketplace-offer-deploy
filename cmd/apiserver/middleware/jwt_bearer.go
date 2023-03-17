package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/security/authentication"
)

const AzureAdJwtKeysUrl = "https://login.microsoftonline.com/common/discovery/v2.0/keys"

// Adds Jwt Bearer authentication to the request
func AddJwtBearer(next http.Handler, config *config.Configuration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validationParameters := getJwtTokenValidationParameters(config)
		isTokenValid := verifyToken(r, validationParameters)

		if !isTokenValid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getJwtTokenValidationParameters(config *config.Configuration) *authentication.JwtTokenValidationParameters {
	keySet, err := authentication.FetchAzureADKeySet(context.Background())

	if err != nil {
		log.Fatal("failed to get Azure AD key set")
	}

	return &authentication.JwtTokenValidationParameters{
		Audience:     config.Azure.ClientId,
		Issuer:       authentication.GetAzureAdIssuer(config.Azure.TenantId),
		IssuerKeySet: keySet,
	}
}

func verifyToken(r *http.Request, parameters *authentication.JwtTokenValidationParameters) bool {
	rawToken := extractToken(r)

	tokenVerifier := authentication.NewJwtTokenVerifier(&rawToken, parameters)
	_, err := tokenVerifier.Verify(r.Context())

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header != "" {
		parts := strings.Split(header, " ")

		if len(parts) == 2 && parts[0] == "Bearer" {
			bearerToken := parts[1]
			return bearerToken
		}
	}
	return ""
}
