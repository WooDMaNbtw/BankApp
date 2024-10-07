package gapi

import (
	"fmt"
	db "github.com/WooDMaNbtw/BankApp/db/sqlc"
	"github.com/WooDMaNbtw/BankApp/pb"
	"github.com/WooDMaNbtw/BankApp/tokens"
	"github.com/WooDMaNbtw/BankApp/utils"
)

// Server service gRPC requests for banking service.
type Server struct {
	pb.UnimplementedBankAppServer
	config     utils.Config
	store      db.Store // interact with database
	tokenMaker tokens.Maker
}

// NewServer creates a new gRPC server and setup routing
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

	return server, nil
}
