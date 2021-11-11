package fiberHandlers

import (
	"errors"
	"stockfyApi/entity"
	"stockfyApi/token"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func MockNewPasetoMaker(symmetricKey string) (token.Maker, error) {

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte("gyqCrSSRpaTVxOtZUxKoBbBOerMfBuVw"),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (
	string, error) {

	// issuedTime := time.Date(2021, time.October, 29, 10, 10, 10, 0, time.UTC)
	// // tokenId, _ := uuid.NewRandom()
	// tokenId, _ := uuid.FromBytes([]byte("TaEHAlSwYIyWfp2s"))
	// fmt.Println("tokenId:", tokenId)
	// payload := token.Payload{
	// 	ID:        tokenId,
	// 	Username:  "TEST_USERNAME",
	// 	IssuedAt:  issuedTime,
	// 	ExpiredAt: issuedTime.Add(time.Minute),
	// }

	// return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return "v2.local.sV-3V-YKYD1n6rIzfUNXj9wKmCILb148U9fqMCYnf7RQo-oi-oEH4bDV" +
		"xP891ZU0N3cgomhNRfOoQo35U3O4iNrURyXipusdOvCE3yVa6I_YaJnkZr43H8_1MLZQ9" +
		"c_x7-TtDYIgjiumj49LkbwmkZdr-B9E8qZ_VaCjbeV52QsnSrxPfxQIzsRInBZtlrrhSo" +
		"6Gk4S3qrG36ufvCH_0SyMg6AxWUnyo8YzGm_w1t1QsurJDUEgz.bnVsbA", nil
}

func (maker *PasetoMaker) VerifyToken(tokenStr string) (*token.Payload, error) {
	// payload := &token.Payload{}

	if tokenStr == "INVALID_TOKEN" {
		return nil, errors.New("invalid token")
	}

	if tokenStr == "EXPIRED_TOKEN" {
		return nil, entity.ErrInvalidApiQueryExpiredToken
	}

	tokenId, _ := uuid.NewRandom()
	return &token.Payload{
		ID:        tokenId,
		Username:  tokenStr,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Minute),
	}, nil
}
