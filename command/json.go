package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/4nth0/golem/config"
	"github.com/4nth0/golem/server"
	jsonServerService "github.com/4nth0/golem/services/json"
	log "github.com/sirupsen/logrus"
)

var RemoteTemplateAddress = "https://raw.githubusercontent.com/4nth0/golem-sample/master/db"

var Version string
var BasePath string = "./.golem"
var TemplatePath string = BasePath + "/templates"
var DatabasePath string = BasePath + "/db"
var ConfigPath string = "./golem.yaml"
var DefaultPort string = "3000"

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

func JsonCmd(ctx context.Context) Command {
	fs := flag.NewFlagSet("golem json", flag.ExitOnError)

	opts := &JsonOpts{}

	fs.StringVar(&opts.path, "path", "", "JSON File")
	fs.StringVar(&opts.port, "port", DefaultPort, "Server port")
	fs.Var(&opts.entities, "entity", "Entity name")
	fs.Var(&opts.templates, "template", "Template name")
	fs.BoolVar(&opts.sync, "sync", true, "FS Sync")

	return Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			fs.Parse(args)
			return Json(ctx, opts)
		},
	}
}

func Json(ctx context.Context, opts *JsonOpts) (err error) {

	if len(opts.entities) == 0 && len(opts.templates) == 0 {
		return errors.New("NO ENTITY PROVIDED")
	}

	defaultServer := server.NewServer(opts.port, nil)

	for _, entity := range opts.entities {
		go initializeEntity(ctx, entity, opts, defaultServer)
	}

	for _, template := range opts.templates {
		if value, ok := DBTemplates[template]; ok {
			err := pullTemplate(template, value)
			if err != nil {
				fmt.Println("Err: ", err)
			}
			initializeEntity(ctx, template, opts, defaultServer)
		}
	}

	fmt.Printf("\nJSON Server has been successfully started and listen on port %s.\n", DefaultPort)

	defaultServer.Listen(ctx)

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

func initializeEntity(ctx context.Context, entity string, opts *JsonOpts, defaultServer *server.Client) {
	entities := map[string]jsonServerService.Entity{
		entity: {
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

	go jsonServerService.LaunchService(ctx, defaultServer, "", service.JSONDBConfig, nil)
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
