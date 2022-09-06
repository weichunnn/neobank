package token

import "time"

// general token maker struct (for different maker)
type Maker interface {
	// create and sign token for a valid username and time duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// check if token is valid or not
	VerifyToken(token string) (*Payload, error)
}
