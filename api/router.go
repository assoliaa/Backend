package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes:= router.Group("/").Use(authMidddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)
	server.router = router
}

func(server *Server)Start(address string)error{ 
	return server.router.Run(address)
}

func errorResponse(err error)gin.H{
	return gin.H{"error": err.Error()}
}