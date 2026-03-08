package server

import (
	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	engine := gin.Default()
	engine.Use(gin.Logger())
	engine.SetTrustedProxies([]string{})
	return engine
}

func NewGroup(server *gin.Engine) *gin.RouterGroup {
	return server.Group("/api/v1")
}
