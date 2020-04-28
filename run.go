package main

import (
	"flag"

	"github.com/4nth0/golem/internal/config"
	"github.com/4nth0/golem/internal/server"
	"github.com/4nth0/golem/internal/services"

	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	configFile string
}

func runCmd() command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &RunOpts{}

	fs.StringVar(&opts.configFile, "config", ConfigPath, "Config File")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return Run(opts)
	}}
}

func Run(opts *RunOpts) (err error) {
	log.Info("Load configuration file")
	cfg := config.LoadConfig(opts.configFile)

	log.Info("Initialize new default server. ", cfg.Port)
	defaultServer := server.NewServer(cfg.Port)

	for _, service := range cfg.Services {
		func(service config.Service) {
			services.Launch(defaultServer, cfg.Vars, service)
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen()
	}

	return nil
}
