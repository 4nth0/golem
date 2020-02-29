package main

import (
	"flag"
	"fmt"

	"github.com/AnthonyCapirchio/golem/internal/config"
	"github.com/AnthonyCapirchio/golem/internal/server"
	jsonServerService "github.com/AnthonyCapirchio/golem/pkg/db/json"
)

type jsonOpts struct {
	path   string
	port   string
	entity stringSlice
	sync   bool
}

// Define a type named "intslice" as a slice of ints
type stringSlice []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

// The second method is Set(value string) error
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func jsonCmd() command {
	fs := flag.NewFlagSet("golem json", flag.ExitOnError)

	opts := &jsonOpts{}

	fs.StringVar(&opts.path, "path", "", "JSON File")
	fs.StringVar(&opts.port, "port", "3000", "Server port")
	fs.Var(&opts.entity, "entity", "Entity name")
	fs.BoolVar(&opts.sync, "sync", true, "FS Sync")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return json(opts)
	}}
}

func json(opts *jsonOpts) (err error) {

	defaultServer := server.NewServer(opts.port)

	for _, entity := range opts.entity {

		entities := map[string]jsonServerService.Entity{
			entity: jsonServerService.Entity{
				DBFile: "./" + entity + ".db.json",
			},
		}

		service := config.Service{
			Port: opts.port,
			JSONDBConfig: jsonServerService.JSONDBConfig{
				Entities: entities,
				Sync:     opts.sync,
			},
		}

		go jsonServerService.LaunchService(defaultServer, "", service.JSONDBConfig)
	}

	defaultServer.Listen()

	return nil
}
