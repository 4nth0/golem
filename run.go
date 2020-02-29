package main

import (
	"flag"

	"github.com/AnthonyCapirchio/golem/internal/config"
	"github.com/AnthonyCapirchio/golem/internal/server"
	"github.com/AnthonyCapirchio/golem/internal/services"
)

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
	defaultServer := server.NewServer(cfg.Port)

	for _, service := range cfg.Services {
		func(service config.Service) {
			services.Launch(defaultServer, service)
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen()
	}

	return nil
}
