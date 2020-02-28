package main

import (
	"fmt"
	"io/ioutil"

	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
	httpService "github.com/AnthonyCapirchio/golem/pkg/server/http"
	"github.com/AnthonyCapirchio/golem/pkg/stats"
	yaml "gopkg.in/yaml.v2"
)

type HttpHandler struct {
	Method string
	Body   string
	Code   int16
}

type GRPCServerConfig struct{}

type Service struct {
	Port         string                         `yaml:"port"`
	Name         string                         `yaml:"name"`
	Type         string                         `yaml:"type"`
	HTTPConfig   httpService.HTTPServerConfig   `yaml:"http_config"`
	JSONDBConfig jsonServerService.JSONDBConfig `yaml:"json_server_config"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func main() {
	s := loadConfig()
	ok := make(chan bool)
	stats := make(chan stats.StatLine)

	for _, service := range s.Services {
		// go httpService.LaunchHttpService(ok, stats, service.Port, service.HTTPConfig)
		func(service Service) {
			if service.Type == "" {
				service.Type = "HTTP"
			}
			switch service.Type {
			case "HTTP":
				go httpService.LaunchHttpService(ok, stats, service.Port, service.HTTPConfig)
			case "JSON_SERVER":
				go jsonServerService.LaunchService(ok, stats, service.Port, service.JSONDBConfig)
			}
		}(service)
	}

	<-ok
}

func launchGRPCService(ok chan<- bool, stats chan<- stats.StatLine, service Service) {
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
