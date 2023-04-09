package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlordtm/pmss/internal/apiserver"
	"github.com/overlordtm/pmss/pkg/pmss"
)

type options struct {
	listenAddr string
}

type Option func(*options)

func WithListenAddr(addr string) Option {
	return func(o *options) {
		o.listenAddr = addr
	}
}

type Server struct {
	options options
	pmss    *pmss.Pmss
	httpSrv *http.Server
}

func New(ctx context.Context, pmss *pmss.Pmss, opts ...Option) *Server {

	ginEngine := apiserver.RegisterHandlersWithOptions(gin.Default(), &handler{pmss}, apiserver.GinServerOptions{
		BaseURL: "/api/v1",
	})

	o := options{
		listenAddr: ":8080",
	}

	for _, opt := range opts {
		opt(&o)
	}

	srv := &Server{
		options: o,
		pmss:    pmss,
	}

	srv.httpSrv = &http.Server{
		Addr:    o.listenAddr,
		Handler: ginEngine,
	}

	return srv
}

func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpSrv.Close()
}
