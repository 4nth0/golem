package services

import (
	"github.com/AnthonyCapirchio/golem/internal/config"
	"github.com/AnthonyCapirchio/golem/internal/server"
	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
	filesServerService "github.com/AnthonyCapirchio/golem/pkg/server/files"
	httpService "github.com/AnthonyCapirchio/golem/pkg/server/http"
	"github.com/gol4ng/logger"
)

// Launch a new service
func Launch(log *logger.Logger, defaultServer *server.Client, globalVars map[string]string, service config.Service) {
	if service.Type == "" {
		service.Type = "HTTP"
	}
	switch service.Type {
	case "HTTP":
		go httpService.LaunchService(log, defaultServer, service.Port, globalVars, service.HTTPConfig)
	case "JSON_SERVER":
		go jsonServerService.LaunchService(defaultServer, service.Port, service.JSONDBConfig)
	case "STATIC":
		go filesServerService.LaunchService(service.Port, service.FilesServerConfig)
	}
}
