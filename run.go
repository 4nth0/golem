package main

import (
	"flag"

	"github.com/AnthonyCapirchio/golem/internal/config"
	"github.com/AnthonyCapirchio/golem/internal/server"
	"github.com/AnthonyCapirchio/golem/internal/services"
	"github.com/gol4ng/logger"
)

type runOpts struct {
	configFile string
}

func runCmd(log *logger.Logger) command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &runOpts{}

	fs.StringVar(&opts.configFile, "config", ConfigPath, "Config File")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return run(log, opts)
	}}
}

func run(log *logger.Logger, opts *runOpts) (err error) {
	log.Info("Load configuration file")
	cfg := config.LoadConfig(opts.configFile)

	log.Info("Initialize new default server. ", logger.String("port", cfg.Port))
	defaultServer := server.NewServer(cfg.Port)

	for _, service := range cfg.Services {
		func(service config.Service) {
			services.Launch(log, defaultServer, service)
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen()
	}

	return nil
}
