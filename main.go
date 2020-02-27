package main

import (
  "fmt"
  "io/ioutil"

  "main/pkg/stats"
  httpService "main/pkg/server/http"
  "gopkg.in/yaml.v2"
)

type HttpHandler struct {
  Method string
  Body string
  Code int16
}

type GRPCServerConfig struct {}

type Service struct {
  Port string `yaml:"port"`
  Name string `yaml:"name"`
  Type string `yaml:"type"`
  HTTPConfig httpService.HTTPServerConfig `yaml:"http_config"`
}

type Config struct {
  Services []Service `yaml:"services"`
}

func main() {

  s := loadConfig()

  ok := make(chan bool)
  stats := make(chan stats.StatLine)

  /* s := Config{
    Services: []Service{
      Service{
        Name: "Move",
        Type: "HTTP",
        Config: httpService.HTTPServerConfig{
          Routes: map[string]httpService.HttpHandler{
            "/ping": httpService.HttpHandler{
              Method: "GET",
              Body: `{"message": "pong!"}`,
              Headers: map[string]string{
                "Content-Type": "application/json",
              },
            },
            "/echo/:message": httpService.HttpHandler{
              Method: "GET",
              Handler: func(w http.ResponseWriter, r *http.Request, params map[string]string) {
                message := params["message"]
                
                w.Write([]byte(message))
              },
              Code: 200,
            },
          },
        },
      },
      Service{
        Name: "Motoko",
        Type: "GRPC",
        Config: GRPCServerConfig{},
      },
    },
  } */

  for _, service := range s.Services {
    go launchService(ok, stats, service)
  }

  <- ok
}

func launchService(ok chan<-bool, stats chan<-stats.StatLine, service Service) {
  switch service.Type {
    case "HTTP":
      httpService.LaunchHttpService(ok, stats, service.Port, service.HTTPConfig)
    case "GRPC":
      launchGRPCService(ok, stats, service)
  }
}

func launchGRPCService(ok chan<-bool, stats chan<-stats.StatLine, service Service) {
  //
}

func loadConfig() *Config {
  t := Config{}

  data, err := ioutil.ReadFile("./golem.yaml")
  if err != nil {
    fmt.Println("Err: ", err)
  }

  err = yaml.Unmarshal(data, &t)
  if err != nil {
    fmt.Println("error: %v", err)
  }

  return &t
}