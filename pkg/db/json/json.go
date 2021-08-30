package json

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/4nth0/golem/pkg/router"
	"github.com/4nth0/golem/pkg/store"
	"github.com/4nth0/golem/server"
)

var empty []byte = []byte("{}")

type JsonServer struct{}

type JSONDBConfig struct {
	Entities map[string]Entity `yaml:"entities"`
	Sync     bool              `yaml:"sync"`
}

type Entity struct {
	DBFile string `yaml:"db_file"`
}

const (
	ContentTypeJSON = "application/json"
	ParamIndexKey   = "index"
)

func LaunchService(ctx context.Context, defaultServer *server.Client, port string, config JSONDBConfig, requests chan server.InboundRequest) {
	var s *server.Client

	if port != "" {
		s = server.NewServer(port, requests)
	} else if defaultServer != nil {
		s = defaultServer
	} else {
		fmt.Println("There is no available server")
		return
	}

	for key, entity := range config.Entities {
		startNewDatabaseServer(key, entity, config.Sync, s.Router)
	}

	if defaultServer == nil {
		defaultServer.Listen(ctx)
	}
}

func startNewDatabaseServer(key string, entity Entity, sync bool, r *router.Router) {
	db := loadDatabaseFromFile(entity.DBFile, sync)

	path := "/" + key
	detailsPath := path + "/:" + ParamIndexKey

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	defer w.Flush()

	r.Get(path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", ContentTypeJSON)

		list, err := json.Marshal(db.List())
		if err != nil {
			fmt.Println("Err: ", err)
		}

		w.Write(list)
	})

	r.Get(detailsPath, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", ContentTypeJSON)

		index, _ := strconv.Atoi(params[ParamIndexKey])
		entry, err := db.GetByIndex(index)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(empty)
			return
		}

		list, err := json.Marshal(entry)
		if err != nil {
			fmt.Println("Err: ", err)
		}

		w.Write(list)
	})

	r.Post(path, func(w http.ResponseWriter, req *http.Request, params map[string]string) {
		decoder := json.NewDecoder(req.Body)
		var t interface{}
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		err = db.Push(t)
		if err != nil {
			fmt.Println("Err: ", err)
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.Delete(detailsPath, func(w http.ResponseWriter, req *http.Request, params map[string]string) {
		index, _ := strconv.Atoi(params["index"])
		err := db.DeleteFromIndex(index)
		if err != nil {
			fmt.Println("Err: ", err)
		}

		w.WriteHeader(http.StatusOK)
	})
}

func loadDatabaseFromFile(path string, sync bool) *store.Database {
	db := store.New(path, sync)
	err := db.Load()
	if err != nil {
		fmt.Println("Err: ", err)
	}

	return db
}
