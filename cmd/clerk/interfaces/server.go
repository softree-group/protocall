package interfaces

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Server struct {
	http http.Server
}

func NewServer(c *ServerConfig, router *mux.Router) *Server {
	return &Server{
		http: http.Server{
			Addr:    fmt.Sprintf("%v:%v", c.Host, c.Port),
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}
