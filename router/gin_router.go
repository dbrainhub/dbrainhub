package router

import (
	"net/http"
	"time"

	"github.com/dbrainhub/dbrainhub/handler"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context) (interface{}, error)

type ginRouter struct {
	server *gin.Engine
}

func (g *ginRouter) Init() {
	g.server = gin.Default()
	g.server.Use(StatMiddleware)
	g.server.Use(gin.Recovery())

	// agent
	g.server.POST("/agent/heartbeat", handlerWapper(handler.Heartbeat))
	g.server.POST("/agent/report", handlerWapper(handler.Report))

	// dbcluster
	g.server.POST("/dbclusters", handlerWapper(handler.CreateDbCluster))
	g.server.GET("/dbclusters", handlerWapper(handler.GetDbClusters))
	// dbcluster memebers
	g.server.GET("/dbclusters/to_assign_members", handlerWapper(handler.GetToAssignedDbClusterMembers))
	g.server.GET("/dbclusters/:clusterId/members", handlerWapper(handler.GetDbClusterMembers))
	g.server.PUT("/dbclusters/:clusterId/members", handlerWapper(handler.AssignDbClusterMembers))
	// tags
	g.server.GET("/tags/all", handlerWapper(handler.GetAllTags))
	g.server.POST("/tags", handlerWapper(handler.AddTag))
	// index
	g.server.POST("/instance/indices", handlerWapper(handler.GetInstanceIndices))
	// dbrainhub output
	g.server.POST("/dbrainhub/slowlogs", handlerWapper(handler.DbRainhubOutput))
	g.server.POST("/dbrainhub/search/instance", handlerWapper(handler.DbRainhubSearchMemberLogsWithCount))
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

func handlerWapper(fun Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := fun(c)
		if err != nil {
			handler.NewApiError(err).RenderJson(c)
		} else {
			c.JSON(http.StatusOK, res)
		}
	}
}
