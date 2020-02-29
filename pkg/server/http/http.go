package http

import (
	"fmt"
	"net/http"

	"github.com/AnthonyCapirchio/golem/internal/server"
	"github.com/AnthonyCapirchio/golem/pkg/stats"
	"github.com/AnthonyCapirchio/golem/pkg/template"
)

type HttpHandler struct {
	Method   string            `yaml:"method"`
	Body     string            `yaml:"body"`
	BodyFile string            `yaml:"body_file"`
	Code     int               `yaml:"code"`
	Headers  map[string]string `yaml:"headers"`
	Handler  *Handler          `yaml:"handler"`
}

type Handler struct {
	Type         string `yaml:"type"`
	Template     string `yaml:"template"`
	TemplateFile string `yaml:"template_file"`
}

type HTTPServerConfig struct {
	Routes map[string]HttpHandler
}

func LaunchService(stats chan<- stats.StatLine, defaultServer *server.ServerClient, port string, config HTTPServerConfig) {

	var s *server.ServerClient

	if port != "" {
		s = server.NewServer(port)
	} else if defaultServer != nil {
		s = defaultServer
	} else {
		fmt.Println("There is no available server")
		return
	}

	for path, route := range config.Routes {
		func(path string, route HttpHandler) {

			if route.Code == 0 {
				route.Code = 200
			}
			if route.Method == "" {
				route.Method = "GET"
			}

			if route.Handler != nil && route.Handler.Type == "template" {
				if route.Handler.TemplateFile != "" {
					route.Handler.Template = template.LoadTemplate(route.Handler.TemplateFile)
				}
			} else if route.Body == "" && route.BodyFile != "" {
				route.Body = template.LoadTemplate(route.BodyFile)
			}

			s.Router.Get(path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {

				if len(route.Headers) > 0 {
					for key, value := range route.Headers {
						w.Header().Add(key, value)
					}
				}
				w.WriteHeader(route.Code)

				if route.Handler != nil {
					switch route.Handler.Type {
					case "template":
						response := template.ExecuteTemplate(route.Handler.Template, params)
						w.Write([]byte(response))
					}
				} else {
					w.Write([]byte(route.Body))
				}

			})
		}(path, route)
	}

	if port != "" {
		s.Listen()
	}
}
