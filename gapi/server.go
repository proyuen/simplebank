package gapi

import (
	"fmt"

	db "github.com/proyuen/simple-bank/db/sqlc"
	"github.com/proyuen/simple-bank/pb"
	"github.com/proyuen/simple-bank/token"
	"github.com/proyuen/simple-bank/util"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

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
