package token

import "time"

type Maker interface { // шобы переключаться между способами создания токена
	CreateToken(username string, duration time.Duration)(string, error)
	
	VerifyToken(token string)(*Payload, error)
}

