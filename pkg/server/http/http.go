package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/AnthonyCapirchio/golem/pkg/stats"
	"github.com/AnthonyCapirchio/t-mux/router"
)

type HttpHandler struct {
	Method   string
	Body     string
	BodyFile string `yaml:"body_file"`
	Code     int
	Headers  map[string]string
	Handler  *Handler
}

type Handler struct {
	Type         string
	Template     string
	TemplateFile string `yaml:"template_file"`
}

type HTTPServerConfig struct {
	Routes map[string]HttpHandler
}

func LaunchHttpService(ok chan<- bool, stats chan<- stats.StatLine, port string, config HTTPServerConfig) {

	s := http.NewServeMux()
	r := router.NewRouter()

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
					route.Handler.Template = loadTemplate(route.Handler.TemplateFile)
				}
			} else if route.Body == "" && route.BodyFile != "" {
				// Use in memory cache
				route.Body = loadTemplate(route.BodyFile)
			}

			r.Get(path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {

				if len(route.Headers) > 0 {
					for key, value := range route.Headers {
						w.Header().Add(key, value)
					}
				}
				w.WriteHeader(route.Code)

				if route.Handler != nil {
					switch route.Handler.Type {
					case "template":
						response := executeTemplate(route.Handler.Template, params)
						w.Write([]byte(response))
					}
				} else {
					w.Write([]byte(route.Body))
				}

			})
		}(path, route)
	}

	s.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		handler, params := r.GetHandler(req.URL.Path, req.Method)
		if handler != nil {
			handler(w, req, params)
		}
	})

	fmt.Println("Starting new server: ", port)

	http.ListenAndServe(":"+port, s)
}

func executeTemplate(template string, params map[string]string) string {
	output := template

	for key, value := range params {
		output = strings.Replace(output, "${"+key+"}", value, -1)
	}

	return output
}

func loadTemplate(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	return string(data)
}
