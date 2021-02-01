package main

import (
	"errors"
	"time"
)

type Authenticator struct {
	tokenizer     Tokenizer
	tokenLifetime time.Duration
}

func NewAuthenticator(
	tokenizer Tokenizer,
	tokenLifetime time.Duration,
) Authenticator {
	return Authenticator{
		tokenizer:     tokenizer,
		tokenLifetime: tokenLifetime,
	}
}

func (a Authenticator) CreateToken(userID int) (string, error) {
	payload := NewTokenPayload(userID, time.Now())
	return a.tokenizer.Encode(payload)
}

func (a Authenticator) GetUserID(token string) (int, error) {
	payload, err := a.tokenizer.Decode(token)
	if err != nil {
		return 0, err
	}

	if payload.UserID < 1 {
		return 0, errors.New("invalid UserID")
	}

	if payload.IssuedAt.Add(a.tokenLifetime).Before(time.Now()) {
		return 0, errors.New("token has expired")
	}

	return payload.UserID, nil
}

type TokenPayload struct {
	UserID   int
	IssuedAt time.Time
}

func NewTokenPayload(
	userID int,
	issuedAt time.Time,
) TokenPayload {
	return TokenPayload{
		UserID:   userID,
		IssuedAt: issuedAt,
	}
}
