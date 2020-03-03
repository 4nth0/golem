package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

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

type stringSlice []string

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
	fs.StringVar(&opts.port, "port", DefaultPort, "Server port")
	fs.Var(&opts.entity, "entity", "Entity name")
	fs.BoolVar(&opts.sync, "sync", true, "FS Sync")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return json(opts)
	}}
}

func json(opts *jsonOpts) (err error) {

	if len(opts.entity) == 0 {
		return errors.New("No entity provided, Please, use at least one entity.")
	}

	defaultServer := server.NewServer(opts.port)

	for _, entity := range opts.entity {

		entities := map[string]jsonServerService.Entity{
			entity: jsonServerService.Entity{
				DBFile: DatabasePath + "/" + entity + ".db.json",
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
		printServiceDetails(entity)
	}

	fmt.Printf("\nJSON Server has been successfully started and listen on port %s.\n", DefaultPort)

	defaultServer.Listen()

	return nil
}

func printServiceDetails(entity string) {
	path := "/" + entity
	detailsPath := path + "/:index"

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 2, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "Method", "Path", "Description")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "------", "----", "-----------")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "GET", path, "Get all resources")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "GET", detailsPath, "Get a specific resource specified by the index")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "POST", path, "Create new resource")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "DELETE", detailsPath, "Delete a specific resource specified by the index")
	fmt.Fprintf(w, "\n\n")
}
