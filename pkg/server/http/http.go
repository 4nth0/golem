package http

import (
	"net/http"

	"github.com/AnthonyCapirchio/golem/internal/server"
	"github.com/AnthonyCapirchio/golem/pkg/template"
	"github.com/gol4ng/logger"
)

// HTTPHandler
type HTTPHandler struct {
	Method   string            `yaml:"method,omitempty"`
	Body     string            `yaml:"body,omitempty"`
	BodyFile string            `yaml:"body_file,omitempty"`
	Code     int               `yaml:"code,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	Handler  *Handler          `yaml:"handler,omitempty"` // Should be removed if not used
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

// LaunchService
func LaunchService(log *logger.Logger, defaultServer *server.Client, port string, globalVars map[string]string, config ServerConfig) {
	var s *server.Client

	log.Info("Launch new HTTP service")

	if port != "" {
		log.Debug("Port provided, create a new server")
		s = server.NewServer(port)
	} else if defaultServer != nil {
		log.Debug("No port provided, use the default server")
		s = defaultServer
	} else {
		log.Info("There is no available server")
		return
	}

	log.Info("Start routes injection")
	for path, route := range config.Routes {
		launch(log, path, route, globalVars, s)
	}

	if port != "" {
		s.Listen()
	}
}

func launch(log *logger.Logger, path string, route HTTPHandler, globalVars map[string]string, s *server.Client) {

	if route.Code == 0 {
		log.Debug("Status code not provided, use default (200)")
		route.Code = 200
	}
	if route.Method == "" {
		log.Debug("HTTP method not provided, use default (GET)")
		route.Method = "GET"
	}

	if route.Body == "" && route.BodyFile != "" {
		log.Info("Use body template file path", logger.String("template_path", route.BodyFile))
		route.Body = template.LoadTemplate(route.BodyFile)
	}

	log.Info("Adding new route", logger.String("method", route.Method), logger.String("path", path))
	s.Router.Add(route.Method, path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {

		log.Info("New inbound request", logger.String("method", route.Method), logger.String("path", path))

		if len(route.Headers) > 0 {
			for key, value := range route.Headers {
				log.Debug("Inject response header", logger.String("key", key), logger.String("value", value))
				w.Header().Add(key, value)
			}
		}

		log.Debug("Use defined status code", logger.Int32("code", int32(route.Code)))
		w.WriteHeader(route.Code)

		response := template.ExecuteTemplate(route.Body, globalVars, params)
		w.Write([]byte(response))

	})
}
