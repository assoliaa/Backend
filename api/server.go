package api

import (
	db "Backend/db/sqlc"
	"Backend/db/utils"
	"Backend/token"
     "fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
	config utils.Config
}

func NewServer(config utils.Config, store db.Store)(*Server, error){
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) //это тоже вынести?
	if err !=nil{
		return nil, fmt.Errorf("could not create token")
	}
	server :=&Server{
		config :config,
		store:store, // и это ?
		tokenMaker:tokenMaker,
	}
	
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency) // и это?
	}
	server.setupRouter()
	return server, nil
}



