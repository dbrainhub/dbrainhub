package router

import (
	"net/http"
	"time"

	"github.com/dbrainhub/dbrainhub/handler"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/gin-gonic/gin"
)

type ginRouter struct {
	server *gin.Engine
}

func (g *ginRouter) Init() {
	g.server = gin.Default()
	g.server.Use(StatMiddleware)
	g.server.Use(gin.Recovery())

	g.server.POST("/hello_world", handler.SayHello)
}

func (g *ginRouter) GetHandler() http.Handler {
	return g.server
}

func StatMiddleware(c *gin.Context) {
	st := time.Now().UnixNano()
	c.Next()

	duration := time.Now().UnixNano() - st
	logger.Infof("%s %s duration: %dus", c.Request.Method, c.Request.URL.Path, duration/10e3)
}
