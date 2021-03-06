package services

import (
	"github.com/4nth0/golem/internal/config"
	jsonServerService "github.com/4nth0/golem/pkg/db/json"
	filesServerService "github.com/4nth0/golem/pkg/server/files"
	httpService "github.com/4nth0/golem/pkg/server/http"
	"github.com/4nth0/golem/server"
)

// Launch a new service
func Launch(defaultServer *server.Client, globalVars map[string]string, service config.Service, requests chan server.InboundRequest) {
	if service.Type == "" {
		service.Type = "HTTP"
	}
	switch service.Type {
	case "HTTP":
		go httpService.LaunchService(defaultServer, service.Port, globalVars, service.HTTPConfig, requests)
	case "JSON_SERVER":
		go jsonServerService.LaunchService(defaultServer, service.Port, service.JSONDBConfig, requests)
	case "STATIC":
		go filesServerService.LaunchService(service.Port, service.FilesServerConfig)
	}
}
