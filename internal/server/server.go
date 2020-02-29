package server

import (
	"net/http"

	"github.com/AnthonyCapirchio/golem/pkg/router"
)

type ServerClient struct {
	Port   string
	Server *http.ServeMux
	Router *router.Router
}

func NewServer(port string) *ServerClient {
	if port == "" {
		return nil
	}

	return &ServerClient{
		Port:   port,
		Server: http.NewServeMux(),
		Router: router.NewRouter(),
	}
}

func (s *ServerClient) Listen() {
	s.Server.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		handler, params := s.Router.GetHandler(req.URL.Path, req.Method)
		if handler != nil {
			handler(w, req, params)
		}
	})

	http.ListenAndServe(":"+s.Port, s.Server)
}
