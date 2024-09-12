package token

import (
	"Backend/db/utils"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T){
	maker, err := NewJWTMaker(utils.RandomString(32))
	assert.NoError(t,err)

	username:= utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
    assert.NotEmpty(t, token)

	payload, err :=maker.VerifyToken(token)
	assert.NoError(t, err)
    assert.NotEmpty(t, payload)

	assert.NotZero(t, payload.ID)
    assert.Equal(t, username, payload.Username)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T){
	maker, err:= NewJWTMaker(utils.RandomString(32))

	assert.NoError(t,err)

	token, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	assert.NoError(t, err)
    assert.NotEmpty(t, token)

	payload, err :=maker.VerifyToken(token)
	assert.Error(t, err)
    assert.EqualError(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)

}

func TestInvalidAlgorithm( t *testing.T){
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
    assert.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	maker, err :=NewJWTMaker(utils.RandomString(32))
	assert.NoError(t, err)


    payload, err = maker.VerifyToken(token)
    assert.Error(t, err)
    assert.EqualError(t, err, ErrInvalidToken.Error())
    assert.Nil(t, payload)
	
}