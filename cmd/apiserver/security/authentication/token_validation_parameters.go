package authentication

import "github.com/lestrrat-go/jwx/jwk"

type JwtTokenValidationParameters struct {
	Issuers       []string
	Audience     string
	IssuerKeySet jwk.Set
}
