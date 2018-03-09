package shared

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const issuer = "test"

type TokenProvider struct{
	secret []byte
}

func baseClaim(expireSeconds int64) jwt.StandardClaims{
	return jwt.StandardClaims{
		ExpiresAt: expireSeconds + time.Now().Unix(),
		Issuer:    issuer,
	}
}

func (tp *TokenProvider) AccountToken(account string) (error, string){
	type MyCustomClaims struct {
		Account  string `json:"account"`
		jwt.StandardClaims
	}

	claims := MyCustomClaims{
		account,
		baseClaim(1500),
	}

	// Create the Claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := jwtToken.SignedString(tp.secret)

	if err != nil {
		return err, ""
	}

	return nil, ss
}

func NewTokenProvider(secret []byte) *TokenProvider{
	return &TokenProvider{secret: secret}
}
