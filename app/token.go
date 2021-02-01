package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Tokenizer interface {
	Encode(payload TokenPayload) (string, error)
	Decode(token string) (TokenPayload, error)
}

var _ Tokenizer = (*JWT)(nil)

type JWT struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

func NewJWT(secret []byte) JWT {
	return JWT{
		secret:        secret,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (j JWT) Encode(payload TokenPayload) (string, error) {
	claims := j.tokenPayloadToClaims(payload)
	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.secret)
}

func (j JWT) Decode(tokenString string) (TokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != j.signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return j.secret, nil
	})

	if err != nil {
		return TokenPayload{}, nil
	}

	if !token.Valid {
		return TokenPayload{}, errors.New("invalid token")
	}

	return j.claimsToTokenPayload(token.Claims)
}

func (JWT) tokenPayloadToClaims(payload TokenPayload) jwt.Claims {
	return jwt.MapClaims{
		"user_id":   payload.UserID,
		"issued_at": payload.IssuedAt,
	}
}

func (JWT) claimsToTokenPayload(claims jwt.Claims) (TokenPayload, error) {
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return TokenPayload{}, errors.New("claims are invalid")
	}

	userID, err := strconv.Atoi(fmt.Sprint(mapClaims["user_id"]))
	if err != nil || userID < 1 {
		return TokenPayload{}, errors.New("invalid user_id")
	}

	issuedAtString, ok := mapClaims["issued_at"].(string)
	if !ok {
		return TokenPayload{}, errors.New("require issued_at")
	}
	issuedAt, err := time.Parse(time.RFC3339, issuedAtString)
	if err != nil {
		return TokenPayload{}, errors.New("invalid issued_at")
	}

	return NewTokenPayload(userID, issuedAt), nil
}
