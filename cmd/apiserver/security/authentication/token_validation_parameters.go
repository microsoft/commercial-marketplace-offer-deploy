package authentication

import "github.com/lestrrat-go/jwx/jwk"

type JwtTokenValidationParameters struct {
	Issuer       string
	Audience     string
	IssuerKeySet jwk.Set
}
