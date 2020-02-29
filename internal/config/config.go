package config

import (
	"fmt"
	"io/ioutil"

	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
	filesServerService "github.com/AnthonyCapirchio/golem/pkg/server/files"
	httpService "github.com/AnthonyCapirchio/golem/pkg/server/http"

	yaml "gopkg.in/yaml.v2"
)

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
