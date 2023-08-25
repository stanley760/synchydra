package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	ServerHTTP *gin.Engine
}

func NewServer(serverHTTP *gin.Engine) *Server {
	return &Server{
		ServerHTTP: serverHTTP,
	}
}
