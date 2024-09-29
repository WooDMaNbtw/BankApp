package api

import (
	db "github.com/WooDMaNbtw/BankApp/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server service HTTP requests for banking service.
type Server struct {
	store  db.Store    // interact with database
	router *gin.Engine // engine for handling requests processing
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default() // setting up default server

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil
		}
	}

	// url patterns initializing
	// accounts apis
	router.GET("/accounts/", server.listAccount)
	router.POST("/accounts/", server.createAccount)
	router.GET("/accounts/:id/", server.getAccount)
	router.PUT("/accounts/:id/", server.updateAccount)
	router.DELETE("/accounts/:id/", server.deleteAccount)

	// transfer apis
	router.POST("/transfers/", server.createTransfer)

	// user apis
	router.POST("/users/", server.createUser)

	// add routes to main router
	server.router = router
	return server
}

// Start run http server on inputted specific HTTP address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
