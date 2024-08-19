package server

import (
	"fmt"
	"net"
	"net/http"

	"example.com/mod/configs"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
}

var server *Server = nil

func NewServer() *Server {
	if server == nil {
		r := chi.NewRouter()
		server = &Server{
			router: r,
			httpServer: &http.Server{
				Addr:    fmt.Sprintf(":d%", configs.GetConfig().Service.APP_PORT),
				Handler: r,
			},
		}
	}

	return server
}

func (s *Server) ListenAddServe() error {
	l, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		panic(err)
	}
	return s.httpServer.Serve(l)
}

func InitServiceEnv() {
	configs.Initialize(".env", ".", ".env")
}

func GetRouter() *chi.Mux {
	return server.router
}

func GetServer() *http.Server {
	return server.httpServer
}
