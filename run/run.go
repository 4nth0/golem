package run

import (
	"context"
	"flag"
	"fmt"

	"github.com/4nth0/golem/internal/command"
	"github.com/4nth0/golem/internal/config"
	"github.com/4nth0/golem/internal/services"
	"github.com/4nth0/golem/server"

	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	ConfigFile   string
	CollectStats bool
}

func RunCmd(ctx context.Context, configPath string) command.Command {
	fs := flag.NewFlagSet("golem run", flag.ExitOnError)

	opts := &RunOpts{}

	fs.StringVar(&opts.ConfigFile, "config", configPath, "Config File")
	fs.BoolVar(&opts.CollectStats, "stats", false, "Collect traffic stats")

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
		go func() {
			for {
				select {
				case <-ctx.Done():
					close(requests)
					return
				case request := <-requests:
					fmt.Println("request: ", request)

				}
			}
		}()
	}

	log.Info("Initialize new default server. ", cfg.Port)
	defaultServer := server.NewServer(cfg.Port, requests)

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
