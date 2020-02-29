package main

import (
	"flag"

	"github.com/AnthonyCapirchio/golem/internal/config"
	"github.com/AnthonyCapirchio/golem/internal/server"
	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
	filesServerService "github.com/AnthonyCapirchio/golem/pkg/server/files"
	httpService "github.com/AnthonyCapirchio/golem/pkg/server/http"
	"github.com/AnthonyCapirchio/golem/pkg/stats"
)

type HttpHandler struct {
	Method string
	Body   string
	Code   int16
}

type runOpts struct {
	configFile string
}

func runCmd() command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &runOpts{}

	fs.StringVar(&opts.configFile, "config", "./golem.yaml", "Config File")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return run(opts)
	}}
}

func run(opts *runOpts) (err error) {
	cfg := config.LoadConfig(opts.configFile)
	stats := make(chan stats.StatLine)
	defaultServer := server.NewServer(cfg.Port)

	for _, service := range cfg.Services {
		func(service config.Service) {
			if service.Type == "" {
				service.Type = "HTTP"
			}
			switch service.Type {
			case "HTTP":
				go httpService.LaunchService(stats, defaultServer, service.Port, service.HTTPConfig)
			case "JSON_SERVER":
				go jsonServerService.LaunchService(stats, defaultServer, service.Port, service.JSONDBConfig)
			case "STATIC":
				go filesServerService.LaunchService(stats, service.Port, service.FilesServerConfig)
			}
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen()
	}

	return nil
}
