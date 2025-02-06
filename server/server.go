package server

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/4nth0/golem/log"
	"github.com/4nth0/golem/router"
)

// Client is the the Server Client
type Client struct {
	Port            string
	Server          *http.ServeMux
	Router          *router.Router
	InboundRequests chan InboundRequest
}

type InboundRequest struct {
	URL     *url.URL            `json:"url"`
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
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
func (s *Client) Listen(ctx context.Context) {
	s.Server.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Info(
			"New inbound request.",
			"method", req.Method,
			"path", req.URL.Path,
			"headers", req.Header,
			"cookies", req.Cookies(),
		)

		s.broadcastInboundRequest(req)

		handler, params, err := s.Router.GetHandler(req.URL.Path, req.Method)
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if handler != nil {
			handler(w, req, params)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: s.Server,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(
				"Unable to start server listening.",
				"err", err,
				"port", s.Port,
			)
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	log.Info("Server stopped")

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Error("server Shutdown Failed", "err", err)
	}
}

func (s *Client) broadcastInboundRequest(req *http.Request) {
	if s.InboundRequests != nil {
		inbound := InboundRequest{
			URL:     req.URL,
			Method:  req.Method,
			Headers: req.Header,
		}

		if req.Body != nil {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Error("Unable to start server listening.", "err", err)
			}
			req.Body.Close()
			inbound.Body = string(body)

			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		s.InboundRequests <- inbound
	}
}
