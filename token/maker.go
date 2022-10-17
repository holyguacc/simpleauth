package token

import "time"

type Maker interface {
	//create token for a username with duration.
	CreateToken(username string, duration time.Duration) (string, error)
	//verify if token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
