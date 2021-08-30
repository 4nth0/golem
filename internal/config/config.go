package config

import (
	"fmt"
	"io/ioutil"

	jsonServerService "github.com/4nth0/golem/pkg/db/json"
	filesServerService "github.com/4nth0/golem/pkg/server/files"
	httpService "github.com/4nth0/golem/pkg/server/http"

	yaml "gopkg.in/yaml.v2"
)

// Service is the Service level configuration struct
type Service struct {
	Port              string                               `yaml:"port,omitempty"`
	Name              string                               `yaml:"name,omitempty"`
	Type              string                               `yaml:"type,omitempty"`
	HTTPConfig        httpService.ServerConfig             `yaml:"http_config,omitempty"`
	JSONDBConfig      jsonServerService.JSONDBConfig       `yaml:"json_server_config,omitempty"`
	FilesServerConfig filesServerService.FilesServerConfig `yaml:"static_server_config,omitempty"`
}

// Config is the rout Config struct
type Config struct {
	path     string
	Vars     map[string]string `yaml:"vars,omitempty"`
	Port     string            `yaml:"port"`
	Services []Service         `yaml:"services"`
}

const FileMode = 0644

// LoadConfig load configuration yaml file content from the specified path
func LoadConfig(path string) (*Config, error) {
	t := Config{
		path: path,
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func InitConfig(path string) *Config {
	cfg := Config{
		path: path,
	}

	return &cfg
}

func (c *Config) SetPort(port string) *Config {
	c.Port = port

	return c
}

func (c Config) Save() error {
	b, _ := yaml.Marshal(c)
	err := ioutil.WriteFile(c.path, b, FileMode)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	return nil
}
