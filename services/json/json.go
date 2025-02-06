package json

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/4nth0/golem/log"
	"github.com/4nth0/golem/router"
	"github.com/4nth0/golem/server"
	"github.com/4nth0/golem/store"
)

type JsonServer struct{}

type JSONDBConfig struct {
	Entities map[string]Entity `yaml:"entities"`
	Sync     bool              `yaml:"sync"`
}

type Entity struct {
	Pagination *Pagination `yaml:"pagination"`
	DBFile     string      `yaml:"db_file"`
}

type Pagination struct {
	Template string `yaml:"template"`
}

const (
	DefaultDBLocalPath = "./.golem/db/%s.json"
)

func LaunchService(ctx context.Context, defaultServer *server.Client, port string, config JSONDBConfig, requests chan server.InboundRequest) {
	var s *server.Client

	log.Info("Launch new JSON Server service")

	if port != "" {
		log.Debug("Port provided, create a new server")
		s = server.NewServer(port, requests)
	} else if defaultServer != nil {
		log.Debug("No port provided, use the default server")
		s = defaultServer
	} else {
		log.Debug("There is no available server")
		return
	}

	for key, entity := range config.Entities {
		log.Info("Create new DB", "key", key)
		mountEntity(key, entity, config.Sync, s.Router)
		printServiceDetails(key)
	}

	if port != "" {
		s.Listen(ctx)
	}
}

func mountEntity(key string, entity Entity, sync bool, r *router.Router) {
	var db *store.Database

	if sync {
		db = loadDatabaseFromFile(key, entity.DBFile, sync)
	} else {
		db = store.New()
	}

	path := "/" + key
	detailsPath := path + "/:" + ParamIndexKey

	r.Get(path, ListHandler(db, key, entity))
	r.Get(detailsPath, GetHandler(db))
	r.Post(path, CreateHandler(db))
	r.Delete(detailsPath, DeleteHandler(db))
}

// @TODO Return error to avoid silent fail
func loadDatabaseFromFile(entity, path string, sync bool) *store.Database {
	if path == "" {
		path = fmt.Sprintf(DefaultDBLocalPath, entity)
	}
	if err := createDBFileIfNotExist(path); err != nil {
		log.Error("Unable to create db file: ", err)
	}

	db := store.New(store.WithLocalFile(path, true))
	err := db.Load()
	if err != nil {
		log.Error("Unable to load database from local file", "err", err)
	}

	return db
}

func createDBFileIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, []byte("[]"), 0644)
	}
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
