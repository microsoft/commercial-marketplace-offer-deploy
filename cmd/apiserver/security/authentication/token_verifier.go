package authentication

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

// Verify JWT tokens
type JwtTokenVerifier interface {
	// Verifies a jwt token
	Verify(ctx context.Context) (*jwt.Token, error)
}

type jwtTokenVerifier struct {
	rawToken   *string
	parameters *JwtTokenValidationParameters
}

// Verifies the JWT token using the provided raw token and the configured validation parameters
func (j *jwtTokenVerifier) Verify(ctx context.Context) (*jwt.Token, error) {
	token, err := parseToken(ctx, j)

	if err != nil {
		return token, err
	}
	return token, verifyClaims(token, j.parameters)
}

// Creates a new JWT token verifier
func NewJwtTokenVerifier(rawToken *string, parameters *JwtTokenValidationParameters) JwtTokenVerifier {
	return &jwtTokenVerifier{rawToken, parameters}
}

// Verifies issuer and audience
func verifyClaims(token *jwt.Token, parameters *JwtTokenValidationParameters) error {
	modmAudience := "api://modm"
	claims, ok := token.Claims.(jwt.MapClaims)
	required := true
	errorMessages := []string{}

	if !ok {
		errorMessages = append(errorMessages, "claims could not be mapped")
	}

	if !claims.VerifyAudience(modmAudience, required) {
		errorMessages = append(errorMessages, fmt.Sprintf("invalid audience %s", parameters.Audience))
	}

	issuerVerified := false
	for _, issuer := range parameters.Issuers {
		if claims.VerifyIssuer(issuer, required) {
			issuerVerified = true
			break
		}
	}
	if !issuerVerified {
		errorMessages = append(errorMessages, fmt.Sprintf("invalid issuer %s", claims["iss"]))
	}

	if len(errorMessages) == 0 {
		return nil
	}
	return utils.NewAggregateError(errorMessages)
}

// Parses the token and verify the signature of the token using the keySet
// on the token verifier
func parseToken(ctx context.Context, j *jwtTokenVerifier) (*jwt.Token, error) {
	token, err := jwt.Parse(*j.rawToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := j.parameters.IssuerKeySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		publickey := &rsa.PublicKey{}
		err := keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}
		return publickey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}
