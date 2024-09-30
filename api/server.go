package api

import (
	"fmt"
	db "github.com/WooDMaNbtw/BankApp/db/sqlc"
	"github.com/WooDMaNbtw/BankApp/tokens"
	"github.com/WooDMaNbtw/BankApp/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server service HTTP requests for banking service.
type Server struct {
	config     utils.Config
	store      db.Store // interact with database
	tokenMaker tokens.Maker
	router     *gin.Engine // engine for handling requests processing
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := tokens.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, fmt.Errorf("validator cannot be registered successfully: %w", err)
		}
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default() // setting up default server
	// url patterns initializing

	// user apis
	router.POST("/users/", server.createUser)
	router.POST("/users/login/", server.loginUser)

	// connecting customized middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// accounts apis
	authRoutes.GET("/accounts/", server.listAccount)
	authRoutes.POST("/accounts/", server.createAccount)
	authRoutes.GET("/accounts/:id/", server.getAccount)
	authRoutes.PUT("/accounts/:id/", server.updateAccount)
	authRoutes.DELETE("/accounts/:id/", server.deleteAccount)

	// transfer apis
	authRoutes.POST("/transfers/", server.createTransfer)

	// add routes to main router
	server.router = router
}

// Start run http server on inputted specific HTTP address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
