package token

import (
	"errors"
	"fmt"
	"stockfyApi/entity"
	"stockfyApi/usecases/utils"
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandString(32))
	require.NoError(t, err)

	username := utils.RandString(10)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(utils.RandString(10), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, entity.ErrInvalidApiQueryExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoSymmetricKeySize(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandString(10))
	require.Equal(t, fmt.Errorf("invalid key size: must be exactly %d characters",
		chacha20poly1305.KeySize), err)
	require.Nil(t, maker)
}

func TestInvalidPasetoTokenDecryption(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandString(32))

	token, err := maker.CreateToken(utils.RandString(10), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	token += "ERROR"
	payload, err := maker.VerifyToken(token)
	fmt.Println(err)
	require.Error(t, err)
	require.EqualError(t, err, errors.New("invalid token authentication").Error())
	require.Nil(t, payload)

}
