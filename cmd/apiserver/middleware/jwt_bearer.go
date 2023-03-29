package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/security/authentication"
)

const AzureAdJwtKeysUrl = "https://login.microsoftonline.com/common/discovery/v2.0/keys"

// Adds Jwt Bearer authentication to the request
func AddJwtBearer(next echo.HandlerFunc, config *config.AppSettings) echo.HandlerFunc {
	return func(c echo.Context) error {
		validationParameters := getJwtTokenValidationParameters(config)
		isTokenValid := verifyToken(c, validationParameters)

		if !isTokenValid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

func getJwtTokenValidationParameters(config *config.AppSettings) *authentication.JwtTokenValidationParameters {
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

func verifyToken(c echo.Context, parameters *authentication.JwtTokenValidationParameters) bool {
	rawToken := extractToken(c)

	tokenVerifier := authentication.NewJwtTokenVerifier(&rawToken, parameters)
	_, err := tokenVerifier.Verify(c.Request().Context())

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func extractToken(c echo.Context) string {
	header := c.Request().Header.Get("Authorization")
	if header != "" {
		parts := strings.Split(header, " ")

		if len(parts) == 2 && parts[0] == "Bearer" {
			bearerToken := parts[1]
			return bearerToken
		}
	}
	return ""
}
