package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

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
		Addr:    srv.options.listenAddr,
		Handler: srv.router(),
	}

	return srv
}

func (s *Server) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/hash/", s.findByHash)
	return mux
}

func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpSrv.Close()
}

func (s *Server) findByHash(w http.ResponseWriter, r *http.Request) {

	hash := r.URL.Path[len("/api/v1/hash/"):]

	result, err := s.pmss.FindByHash(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	payload := HashResponse{
		Hash:        hash,
		HashVariant: result.HashVariant,
		Status:      HashStatusSafe,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
