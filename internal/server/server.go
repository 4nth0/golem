package server

import (
	"fmt"
	"net/http"

	"github.com/AnthonyCapirchio/golem/pkg/router"
	log "github.com/sirupsen/logrus"
)

// Client is the the Server Client
type Client struct {
	Port   string
	Server *http.ServeMux
	Router *router.Router
}

// NewServer create a new Server instance that contains a new Mux and a new Router
func NewServer(port string) *Client {
	if port == "" {
		return nil
	}

	return &Client{
		Port:   port,
		Server: http.NewServeMux(),
		Router: router.NewRouter(),
	}
}

// Listen start listening
func (s *Client) Listen() {
	s.Server.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.WithFields(
			log.Fields{
				"method": req.Method,
				"path":   req.URL.Path,
			}).Info("New inbound request.")
		handler, params := s.Router.GetHandler(req.URL.Path, req.Method)
		if handler != nil {
			handler(w, req, params)
		}
	})

	err := http.ListenAndServe(":"+s.Port, s.Server)
	if err != nil {
		fmt.Println("Err: ", err)
	}
}
