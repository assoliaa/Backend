package token

import (
	"fmt"
	"time"
	"github.com/o1egl/paseto"
)

const (
	// Размер ключа для chacha20poly1305
	symmetricKeySize = 32
)

// PasetoMaker используется для создания и проверки PASETO токенов
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey  []byte
}


func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != symmetricKeySize {
		return nil, fmt.Errorf("invalid key size: must be %d bytes", symmetricKeySize)
	}
	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}


func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err !=nil{
		return "", err
	}
	
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}


func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload:= &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
    if err !=nil{
		return nil, ErrInvalidToken
	}
	err =payload.Valid()
	if err !=nil{
		return nil, err
	}
	return payload, nil
}
