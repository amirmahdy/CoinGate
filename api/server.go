package api

import (
	"db"
	"fmt"
	"token"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("coin", validCoin)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/login", server.loginUser)
	// router.POST("/token/refresh", server.handleTokenRefresh)

	// authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// authRouters.POST("/accounts", server.createAccount)
	// authRouters.GET("/accounts/:username", server.getAccount)

	// authRouters.POST("/transfer", server.createTransfer)

	// authRouters.GET("/coin", server.listCoin)
	// authRouters.POST("/coin", server.createCoin)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
