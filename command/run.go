package command

import (
	"context"
	"flag"

	"github.com/4nth0/golem/config"
	"github.com/4nth0/golem/server"
	"github.com/4nth0/golem/services"
	"github.com/4nth0/golem/stats"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	ConfigFile       string
	CollectStats     bool
	StatsDestination string
	StatsDriver      string
	Debug            bool
}

func RunCmd(ctx context.Context, configPath string) Command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &RunOpts{}

	fs.StringVar(&opts.ConfigFile, "config", configPath, "Config File")

	fs.BoolVar(&opts.CollectStats, "stats", false, "Collect traffic stats")
	fs.StringVar(&opts.StatsDestination, "stats-dest", "./stats.log", "Collected traffic destination")
	fs.StringVar(&opts.StatsDriver, "stats-driver", "fs", "Collected traffic destination")

	return Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			err := fs.Parse(args)
			if err != nil {
				return err
			}
			return Run(ctx, opts, nil)
		},
	}
}

func Run(ctx context.Context, opts *RunOpts, requests chan server.InboundRequest) error {
	log.Info("Load configuration file")
	cfg, err := config.LoadConfig(opts.ConfigFile)
	if err != nil {
		return err
	}

	if opts.CollectStats {
		if requests == nil {
			requests = make(chan server.InboundRequest)
		}
		col, err := stats.NewCollector(opts.StatsDriver, opts.StatsDestination)
		if err != nil {
			return err
		}

		go col.Collect(ctx, requests)
	}

	log.Info("Initialize new default server. ", cfg.Port)
	defaultServer := server.NewServer(cfg.Port, requests)
	defaultServer.Server.Handle("/metrics", promhttp.Handler())

	for _, service := range cfg.Services {
		func(service config.Service) {
			services.Launch(ctx, defaultServer, cfg.Vars, service, requests)
		}(service)
	}

	if defaultServer != nil {
		defaultServer.Listen(ctx)
	}

	return nil
}
