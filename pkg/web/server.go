package web

import (
	"fmt"
	"net/http"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Server struct {
	http http.Server
}

func NewServer(c *ServerConfig, mux *http.ServeMux) *Server {
	return &Server{
		http: http.Server{
			Addr:    fmt.Sprintf("%v:%v", c.Host, c.Port),
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}
