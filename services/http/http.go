package http

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/4nth0/golem/server"
	log "github.com/sirupsen/logrus"
)

// HTTPHandler
type HTTPHandler struct {
	Method    string                 `yaml:"method,omitempty"`
	Methods   map[string]HTTPHandler `yaml:"methods,omitempty"`
	Body      string                 `yaml:"body,omitempty"`
	Bodies    []string               `yaml:"bodies,omitempty"`
	BodyFile  string                 `yaml:"body_file,omitempty"`
	BodyFiles []string               `yaml:"body_files,omitempty"`
	Code      int                    `yaml:"code,omitempty"`
	Headers   map[string]string      `yaml:"headers,omitempty"`
	Handler   *Handler               `yaml:"handler,omitempty"` // Should be removed if not used
	Latency   time.Duration          `yaml:"latency,omitempty"` // Should be removed if not used
}

// Handler
type Handler struct {
	Type         string `yaml:"type"`
	Template     string `yaml:"template"`
	TemplateFile string `yaml:"template_file"`
}

// ServerConfig
type ServerConfig struct {
	Routes map[string]HTTPHandler
}

var (
	DefaultMethod     string = "GET"
	DefaultStatusCode int    = http.StatusOK
)

// LaunchService
func LaunchService(ctx context.Context, defaultServer *server.Client, port string, globalVars map[string]string, config ServerConfig, requests chan server.InboundRequest) {
	var s *server.Client

	log.Info("Launch new HTTP service")

	if port != "" {
		log.Debug("Port provided, create a new server")
		s = server.NewServer(port, requests)
	} else if defaultServer != nil {
		log.Debug("No port provided, use the default server")
		s = defaultServer
	} else {
		log.Info("There is no available server")
		return
	}

	log.Info("Start routes injection")
	for path, route := range config.Routes {
		if len(route.Methods) > 0 {
			for method, route := range route.Methods {
				route.Method = method
				launch(path, route, globalVars, s)
			}
		} else {
			launch(path, route, globalVars, s)
		}
	}

	if port != "" {
		s.Listen(ctx)
	}
}

func launch(path string, route HTTPHandler, globalVars map[string]string, s *server.Client) {

	route = normalizeRouteConfiguration(route)

	log.WithFields(
		log.Fields{
			"method": route.Method,
			"path":   path,
		}).Info("Adding new route")

	s.Router.Add(route.Method, path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		log.WithFields(
			log.Fields{
				"method": route.Method,
				"path":   path,
				"status": route.Code,
			}).Info("New inbound request.")

		if len(route.Headers) > 0 {
			for key, value := range route.Headers {
				w.Header().Add(key, value)
			}
		}
		if route.Latency > 0 {
			time.Sleep(route.Latency)
		}

		w.WriteHeader(route.Code)
		w.Write([]byte(generateBodyResponse(route, globalVars, params)))

	})
}

func normalizeRouteConfiguration(route HTTPHandler) HTTPHandler {
	if route.Code == 0 {
		log.WithFields(
			log.Fields{
				"code": DefaultStatusCode,
			}).Debug("Status code not provided, use default.")

		route.Code = DefaultStatusCode
	}
	if route.Method == "" {
		log.WithFields(
			log.Fields{
				"method": DefaultMethod,
			}).Debug("HTTP method not provided, use default.")

		route.Method = DefaultMethod
	}

	if route.BodyFile != "" {
		log.WithFields(
			log.Fields{
				"path": route.BodyFile,
			}).Debug("Use body template file.")

		result, err := LoadTemplate(route.BodyFile)
		if err != nil {
			log.WithFields(
				log.Fields{
					"path": route.BodyFile,
				}).Error("Unable to load template file")
		} else {
			route.Body = result
		}
	}

	if len(route.BodyFiles) > 0 {
		log.WithFields(
			log.Fields{
				"paths": route.BodyFiles,
			}).Debug("Use body template file.")

		for _, path := range route.BodyFiles {
			result, err := LoadTemplate(path)
			if err != nil {
				log.WithFields(
					log.Fields{
						"path": path,
					}).Error("Unable to load template file")
			} else {
				route.Bodies = append(route.Bodies, result)
			}
		}
	}

	return route
}

func generateBodyResponse(route HTTPHandler, globalVars map[string]string, params map[string]string) string {
	var body string

	if len(route.Bodies) > 0 {
		body = route.Bodies[rand.Intn(len(route.Bodies))]
	} else {
		body = route.Body
	}

	return ExecuteTemplate(body, globalVars, params)
}
