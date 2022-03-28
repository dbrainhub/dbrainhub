package router

type HTTPRouter interface {
	Run() error
}

func NewDefaultRouter(addr string) HTTPRouter {
	var router defaultRouter

	router.Init(addr)
	return &router
}
