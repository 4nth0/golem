package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/AnthonyCapirchio/golem/internal/server"
	"github.com/AnthonyCapirchio/golem/pkg/router"
	"github.com/AnthonyCapirchio/golem/pkg/store"
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

func LaunchService(defaultServer *server.Client, port string, config JSONDBConfig) {
	var s *server.Client

	if port != "" {
		s = server.NewServer(port)
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
		defaultServer.Listen()
	}
}

func startNewDatabaseServer(key string, entity Entity, sync bool, r *router.Router) {
	db := loadDatabaseFromFile(entity.DBFile, sync)

	path := "/" + key
	detailsPath := path + "/:index"

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	defer w.Flush()

	r.Get(path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", "application/json")

		list, err := json.Marshal(db.List())
		if err != nil {
			fmt.Println("Err: ", err)
		}

		w.Write(list)
	})

	r.Get(detailsPath, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", "application/json")

		index, _ := strconv.Atoi(params["index"])
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
	db := store.Database{
		FilePath: path,
		Sync:     sync,
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			db.InitDefault()
		}
	} else {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("Err: ", err)
		}

		err = json.Unmarshal(data, &db)
		if err != nil {
			fmt.Println("Err: ", err)
		}
	}

	return &db
}
