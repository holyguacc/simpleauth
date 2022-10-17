package api

import (
	"fmt"
	db "simpleauth/db/sqlc"
	"simpleauth/token"
	"simpleauth/util"

	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store, mailer *mail.Msg) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()
	user := router.Group("/user")
	{
		user.POST("/register", server.createUser)
		user.POST("/login", server.LoginUser)
	}
	authorized := router.Group("/auth").Use(authMiddleware(server.tokenMaker))
	{
		authorized.POST("/posts", server.CreatePost)
	}
	router.GET("/posts", server.ListPosts)

	verify := router.Group("/verify")
	{
		verify.GET("/send", server.verifyUser)
		verify.POST("/get", server.sendVerificationCode)
	}

	router.POST("/forgot", server.SendResetPassword)
	router.POST("/resetpassword", server.ResetAndChangePassword)

	router.Run()
	return server, nil
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
