package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(tokenStr string) (*Payload, error)
}

// type NewTokenMaker struct {
// 	TokenMaker func(symmetricKey string) (Maker, error)
// }

// type NewPasetoMakerInterface interface {
// 	NewPasetoMaker(symmetricKey string) (Maker, error)
// }
