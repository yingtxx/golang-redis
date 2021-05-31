package http

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	*http.Server
}

type ServerConf struct {
	Addr string
}

func NewServer(cf *ServerConf) (svr *Server) {
	return &Server{
		Server: &http.Server{
			Addr: cf.Addr,
		},
	}
}

func (svr *Server) SetHandler(handler http.Handler) {
	svr.Handler = handler
}

func (svr *Server) Run() (err error) {
	if err = svr.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		err = fmt.Errorf("http.Server.ListenAndServe: %w", err)
		return
	}

	err = nil

	return
}

func (svr *Server) Close() (err error) {
	if err = svr.Shutdown(context.TODO()); err != nil {
		err = fmt.Errorf("server shutdown: %w", err)
		return
	}

	return
}
