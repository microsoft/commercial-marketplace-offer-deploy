package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

const AzureAdJwtKeysUrl = "https://login.microsoftonline.com/common/discovery/v2.0/keys"

// Adds Jwt Bearer authentication to the request
func AddJwtBearer(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		token, err := verifyToken(r.Context(), &tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			// TODO: convert to debug call
			log.Printf("  JWT Bearer Token Valid for Issuer=%s", token.Claims.(jwt.MapClaims)["aud"])
		}
		next.ServeHTTP(w, r)
	})
}

func verifyToken(ctx context.Context, tokenString *string) (*jwt.Token, error) {
	keySet, err := jwk.Fetch(ctx, AzureAdJwtKeysUrl)

	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
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
		return nil, err
	}

	return token, nil
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
