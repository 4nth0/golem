package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/4nth0/golem/internal/command"
	"github.com/4nth0/golem/internal/config"
	jsonServerService "github.com/4nth0/golem/pkg/db/json"
	"github.com/4nth0/golem/server"
	log "github.com/sirupsen/logrus"
)

var RemoteTemplateAddress = "https://raw.githubusercontent.com/4nth0/golem-sample/master/db"

var DBTemplates = map[string]string{
	"users": RemoteTemplateAddress + "/users/db.json",
}

type JsonOpts struct {
	path      string
	port      string
	entities  stringSlice
	templates stringSlice
	sync      bool
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

func jsonCmd() command.Command {
	fs := flag.NewFlagSet("golem json", flag.ExitOnError)

	opts := &JsonOpts{}

	fs.StringVar(&opts.path, "path", "", "JSON File")
	fs.StringVar(&opts.port, "port", DefaultPort, "Server port")
	fs.Var(&opts.entities, "entity", "Entity name")
	fs.Var(&opts.templates, "template", "Template name")
	fs.BoolVar(&opts.sync, "sync", true, "FS Sync")

	return command.Command{fs, func(args []string) error {
		fs.Parse(args)
		return Json(opts)
	}}
}

func Json(opts *JsonOpts) (err error) {

	if len(opts.entities) == 0 && len(opts.templates) == 0 {
		return errors.New("No entity provided, Please, use at least one entity.")
	}

	defaultServer := server.NewServer(opts.port, nil)

	for _, entity := range opts.entities {
		go initializeEntity(entity, opts, defaultServer)
	}

	for _, template := range opts.templates {
		if value, ok := DBTemplates[template]; ok != false {
			err := pullTemplate(template, value)
			if err != nil {
				fmt.Println("Err: ", err)
			}
			initializeEntity(template, opts, defaultServer)
		}
	}

	fmt.Printf("\nJSON Server has been successfully started and listen on port %s.\n", DefaultPort)

	defaultServer.Listen()

	return nil
}

func pullTemplate(entity, path string) error {

	log.WithFields(
		log.Fields{
			"entity": entity,
			"path":   path,
		}).Info("Pull DB template.")

	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(DatabasePath + "/" + entity + ".db.json")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func initializeEntity(entity string, opts *JsonOpts, defaultServer *server.Client) {
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

	go jsonServerService.LaunchService(defaultServer, "", service.JSONDBConfig, nil)
	printServiceDetails(entity)
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
