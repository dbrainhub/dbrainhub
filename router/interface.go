package router

import "net/http"

func NewDefaultHandler() http.Handler {
	var router ginRouter

	router.Init()
	return router.GetHandler()
}
