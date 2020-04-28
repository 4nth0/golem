package server

import (
	"net/http"
	"net/url"

	"github.com/4nth0/golem/pkg/router"
	log "github.com/sirupsen/logrus"
)

// Client is the the Server Client
type Client struct {
	Port            string
	Server          *http.ServeMux
	Router          *router.Router
	InboundRequests chan InboundRequest
}

type InboundRequest struct {
	URL    *url.URL
	Method string
}

// NewServer create a new Server instance that contains a new Mux and a new Router
func NewServer(port string, requests chan InboundRequest) *Client {
	if port == "" {
		return nil
	}

	return &Client{
		Port:            port,
		Server:          http.NewServeMux(),
		Router:          router.NewRouter(),
		InboundRequests: requests,
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
			if s.InboundRequests != nil {
				s.InboundRequests <- InboundRequest{
					URL:    req.URL,
					Method: req.Method,
				}
			}
			handler(w, req, params)
		}
	})

	err := http.ListenAndServe(":"+s.Port, s.Server)
	if err != nil {
		log.WithFields(
			log.Fields{
				"err":  err,
				"port": s.Port,
			}).Error("Unable to start server listening.")
	}
}
