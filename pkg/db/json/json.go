package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/AnthonyCapirchio/golem/internal/server"
	"github.com/AnthonyCapirchio/golem/pkg/stats"
	"github.com/AnthonyCapirchio/golem/pkg/store"
	"github.com/AnthonyCapirchio/t-mux/router"
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

func LaunchService(stats chan<- stats.StatLine, defaultServer *server.ServerClient, port string, config JSONDBConfig) {

	var s *server.ServerClient

	if port != "" {
		s = server.NewServer(port)
	} else if defaultServer != nil {
		s = defaultServer
	} else {
		fmt.Println("There is no available server")
		return
	}

	for key, entity := range config.Entities {
		go startNewDatabaseServer(key, entity, config.Sync, s.Router)
	}

	s.Listen()
}

func startNewDatabaseServer(key string, entity Entity, sync bool, r *router.Router) {
	db := loadDatabaseFromFile(entity.DBFile, sync)

	path := "/" + key
	detailsPath := path + "/:index"

	r.Get(path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		list, err := json.Marshal(db.List())
		if err != nil {
			fmt.Println("Err: ", err)
		}

		w.Write(list)
	})

	r.Tree.AddNode(detailsPath, "GET", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	r.Tree.AddNode(path, "POST", func(w http.ResponseWriter, req *http.Request, params map[string]string) {
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

	r.Tree.AddNode(detailsPath, "DELETE", func(w http.ResponseWriter, req *http.Request, params map[string]string) {
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
