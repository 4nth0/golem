package run

import (
	"context"
	"flag"

	"github.com/4nth0/golem/internal/command"
	"github.com/4nth0/golem/internal/config"
	"github.com/4nth0/golem/internal/services"
	"github.com/4nth0/golem/pkg/stats"
	"github.com/4nth0/golem/server"

	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	ConfigFile       string
	CollectStats     bool
	StatsDestination string
	StatsDriver      string
}

func RunCmd(ctx context.Context, configPath string) command.Command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &RunOpts{}

	fs.StringVar(&opts.ConfigFile, "config", configPath, "Config File")

	fs.BoolVar(&opts.CollectStats, "stats", false, "Collect traffic stats")
	fs.StringVar(&opts.StatsDestination, "stats-dest", "./stats.log", "Collected traffic destination")
	fs.StringVar(&opts.StatsDriver, "stats-driver", "fs", "Collected traffic destination")

	return command.Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			fs.Parse(args)
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

	defaultServer.Server.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	defaultServer.Server.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	defaultServer.Server.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	defaultServer.Server.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	defaultServer.Server.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	defaultServer.Server.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))

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
