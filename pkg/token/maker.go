package token

import "time"

type Maker interface {
	CreateToken(username string, userID int64, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
