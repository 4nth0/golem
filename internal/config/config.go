package config

import (
	"fmt"
	"io/ioutil"

	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
	filesServerService "github.com/AnthonyCapirchio/golem/pkg/server/files"
	httpService "github.com/AnthonyCapirchio/golem/pkg/server/http"

	yaml "gopkg.in/yaml.v2"
)

// Service is the Service level configuration struct
type Service struct {
	Port              string                               `yaml:"port"`
	Name              string                               `yaml:"name"`
	Type              string                               `yaml:"type"`
	HTTPConfig        httpService.ServerConfig             `yaml:"http_config"`
	JSONDBConfig      jsonServerService.JSONDBConfig       `yaml:"json_server_config"`
	FilesServerConfig filesServerService.FilesServerConfig `yaml:"static_server_config"`
}

// Config is the rout Config struct
type Config struct {
	Port     string    `yaml:"port"`
	Services []Service `yaml:"services"`
}

// LoadConfig load configuration yaml file content from the specified path
func LoadConfig(path string) *Config {
	t := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		fmt.Println("error: %v", err)
	}

	return &t
}
