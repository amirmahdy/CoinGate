package token

import (
	"errors"
	"time"
)

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

var (
	ErrInvalidSecretKey = errors.New("invalid secret key")
	ErrExpiredToken     = errors.New("token is expired")
	ErrorInvalidToken   = errors.New("token is invalid")
)
