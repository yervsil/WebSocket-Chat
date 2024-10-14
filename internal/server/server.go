package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yervsil/auth_service/internal/configs"
)

type HttpServer struct {
	Server http.Server
}

func New(config *configs.Config, router *mux.Router) *HttpServer {
	return &HttpServer{
		Server: http.Server{
					Addr: config.HttpServer.Port,
					ReadTimeout: config.ReadTimeout,
					WriteTimeout: config.WriteTimeout,
					Handler: router,
				},
	}
}

func(s *HttpServer) Run() error{
	return s.Server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}