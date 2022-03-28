package router

import (
	"context"
	"fmt"
	nativehttp "net/http"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type defaultRouter struct {
	server *http.Server
	app    *kratos.App
}

func (d *defaultRouter) Init(addr string) {
	d.server = http.NewServer(
		http.Address(addr),
		http.Middleware(
			// add service filter
			StatMiddleware,
			RecoveryMiddleware(),
		),
		http.ErrorEncoder(ErrorsHandler),
		http.ResponseEncoder(DefaultResponseEncoder),
	)

	api.RegisterHttpServiceHTTPServer(d.server, &controller.DefaultController{})

	d.app = kratos.New(
		kratos.Name("dbrainhub"),
		kratos.Server(
			d.server,
		),
	)
}

func (d *defaultRouter) Run() error {
	return d.app.Run()
}

func StatMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		st := time.Now().UnixNano()

		reply, err = handler(ctx, req)

		duration := time.Now().UnixNano() - st
		info := ctx.(http.Context)
		logger.Infof("%s %s duration: %dus", info.Request().Method, info.Request().URL, duration/1e3)
		return
	}
}

func RecoveryMiddleware() middleware.Middleware {
	return recovery.Recovery(recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
		logger.Errorf("panic happends, req: %v, err: %v", req, err)
		return fmt.Errorf("panic error")
	}))
}

func ErrorsHandlerMiddlerware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		reply, err = handler(ctx, req)

		if err == nil {
			return
		}
		if _, ok := err.(*errors.ErrInfo); !ok {
			err = errors.NewErrInfo(-1, "unknown error", err.Error())
		}

		return err.(*errors.ErrInfo).Error(), nil
	}
}

func ErrorsHandler(writer nativehttp.ResponseWriter, req *nativehttp.Request, err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*errors.ErrInfo); !ok {
		writer.WriteHeader(nativehttp.StatusInternalServerError)
		err = errors.NewUnknownError()
	} else {
		writer.WriteHeader(nativehttp.StatusOK)
	}

	_, _ = writer.Write([]byte(err.Error()))
}

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(w nativehttp.ResponseWriter, r *nativehttp.Request, v interface{}) error {
	succResp := errors.NewSuccessResp(v)

	w.WriteHeader(nativehttp.StatusOK)
	w.Header().Set("Content-type", "application/json")
	_, err := w.Write([]byte(succResp.String()))
	return err
}
