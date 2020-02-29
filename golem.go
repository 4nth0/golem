package main

import (
	"fmt"
	"io/ioutil"

	"github.com/AnthonyCapirchio/golem/internal/server"
	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
	filesServerService "github.com/AnthonyCapirchio/golem/pkg/server/files"
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
	Port              string                               `yaml:"port"`
	Name              string                               `yaml:"name"`
	Type              string                               `yaml:"type"`
	HTTPConfig        httpService.HTTPServerConfig         `yaml:"http_config"`
	JSONDBConfig      jsonServerService.JSONDBConfig       `yaml:"json_server_config"`
	FilesServerConfig filesServerService.FilesServerConfig `yaml:"static_server_config"`
}

type Config struct {
	Port     string    `yaml:"port"`
	Services []Service `yaml:"services"`
}

func main() {
	config := loadConfig()
	ok := make(chan bool)
	stats := make(chan stats.StatLine)
	defaultServer := server.NewServer(config.Port)

	for _, service := range config.Services {
		// go httpService.LaunchHttpService(ok, stats, service.Port, service.HTTPConfig)
		func(service Service) {
			if service.Type == "" {
				service.Type = "HTTP"
			}
			switch service.Type {
			case "HTTP":
				go httpService.LaunchService(ok, stats, defaultServer, service.Port, service.HTTPConfig)
			case "JSON_SERVER":
				go jsonServerService.LaunchService(ok, stats, defaultServer, service.Port, service.JSONDBConfig)
			case "STATIC":
				go filesServerService.LaunchService(ok, stats, service.Port, service.FilesServerConfig)
			}
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen()
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
