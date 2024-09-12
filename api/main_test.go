package api

import (
	db "Backend/db/sqlc"
	"Backend/db/utils"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func newTestServer(t *testing.T, store db.Store)*Server{ // here
	config:=utils.Config{
		TokenSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	assert.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
    gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}