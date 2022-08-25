package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/weichunnn/neobank/db/sqlc"
)

// serve http requests for banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// create server and setup http routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)

	server.router = router
	return server
}

// runs http server on a address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
