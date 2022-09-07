package gapi

import (
	"fmt"

	db "github.com/weichunnn/neobank/db/sqlc"
	"github.com/weichunnn/neobank/pb"
	"github.com/weichunnn/neobank/token"
	"github.com/weichunnn/neobank/util"
)

// serve gRPC requests for banking service
type Server struct {
	pb.UnimplementedNeoBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// create a new gRPC
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
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
