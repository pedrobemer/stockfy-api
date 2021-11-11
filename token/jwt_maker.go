package token

import (
	"errors"
	"fmt"
	"stockfyApi/entity"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters",
			minSecretKeySize)

	}

	return &JWTMaker{secretKey: secretKey}, nil

}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (
	string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("INVALID_TOKEN_SIGNING_METHOD")
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		fmt.Println(verr.Inner)
		if ok && errors.Is(verr.Inner, entity.ErrInvalidApiQueryExpiredToken) {
			return nil, entity.ErrInvalidApiQueryExpiredToken
		}

		return nil, entity.ErrInvalidApiQueryInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, entity.ErrInvalidApiQueryInvalidToken
	}

	return payload, nil
}
