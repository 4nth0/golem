package run

import (
	"flag"

	"github.com/4nth0/golem/internal/command"
	"github.com/4nth0/golem/internal/config"
	"github.com/4nth0/golem/internal/services"
	"github.com/4nth0/golem/server"

	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	ConfigFile string
}

func RunCmd(configPath string) command.Command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &RunOpts{}

	fs.StringVar(&opts.ConfigFile, "config", configPath, "Config File")

	return command.Command{fs, func(args []string) error {
		fs.Parse(args)
		return Run(opts, nil)
	}}
}

func Run(opts *RunOpts, requests chan server.InboundRequest) error {
	log.Info("Load configuration file")
	cfg := config.LoadConfig(opts.ConfigFile)

	log.Info("Initialize new default server. ", cfg.Port)
	defaultServer := server.NewServer(cfg.Port, requests)

	for _, service := range cfg.Services {
		func(service config.Service) {
			services.Launch(defaultServer, cfg.Vars, service, requests)
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen()
	}

	return nil
}
