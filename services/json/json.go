package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/4nth0/golem/router"
	"github.com/4nth0/golem/server"
	"github.com/4nth0/golem/store"

	log "github.com/sirupsen/logrus"
)

var empty []byte = []byte("{}")

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
	ContentTypeJSON         = "application/json"
	ParamIndexKey           = "index"
	DefaultPaginationtLimit = 10
	DefaultPaginationtPage  = 0
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
		fmt.Println("There is no available server")
		return
	}

	for key, entity := range config.Entities {
		log.Info("Create new DB for: ", key)
		startNewDatabaseServer(key, entity, config.Sync, s.Router)
		printServiceDetails(key)
	}

	if port != "" {
		s.Listen(ctx)
	}
}

func startNewDatabaseServer(key string, entity Entity, sync bool, r *router.Router) {
	var db *store.Database

	if sync {
		db = loadDatabaseFromFile(key, entity.DBFile, sync)
	} else {
		db = store.New()
	}

	path := "/" + key
	detailsPath := path + "/:" + ParamIndexKey

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	defer w.Flush()

	r.Get(path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", ContentTypeJSON)

		var list []byte
		var err error

		if entity.Pagination != nil {
			list, err = renderPaginatedList(entity, r, db)
			if err != nil {
				fmt.Println("Err: ", err)
			}
		} else {
			list, err = json.Marshal(db.List())
			if err != nil {
				fmt.Println("Err: ", err)
			}
		}

		_, err = w.Write(list)
		if err != nil {
			log.WithFields(
				log.Fields{
					"err": err,
				}).Error("Unable to write response.")
		}
	})

	r.Get(detailsPath, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", ContentTypeJSON)

		index, _ := strconv.Atoi(params[ParamIndexKey])
		entry, err := db.GetByIndex(index)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write(empty)
			log.WithFields(
				log.Fields{
					"err": err,
				}).Error("Unable to write response.")
			return
		}

		list, err := json.Marshal(entry)
		if err != nil {
			fmt.Println("Err: ", err)
		}

		_, err = w.Write(list)
		if err != nil {
			log.WithFields(
				log.Fields{
					"err": err,
				}).Error("Unable to write response.")
		}
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

func loadDatabaseFromFile(entity, path string, sync bool) *store.Database {
	if path == "" {
		path = fmt.Sprintf("./.golem/db/%s.json", entity)
	}
	if err := createDBFileIfNotExist(path); err != nil {
		log.Error("Unable to create db file: ", err)
	}

	db := store.New(store.WithLocalFileSync(path))
	err := db.Load()
	if err != nil {
		fmt.Println("Err: ", err)
	}

	return db
}

func createDBFileIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, []byte("[]"), 0644)
	}
	return nil
}

func renderPaginatedList(entity Entity, r *http.Request, db *store.Database) ([]byte, error) {
	var err error
	limit := DefaultPaginationtLimit
	if r.URL.Query().Get("limit") != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 1 {
			limit = DefaultPaginationtLimit
		}
	}

	page := DefaultPaginationtPage
	if r.URL.Query().Get("page") != "" {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 0 {
			page = 0
		}
	}

	entries := db.PaginatedList(page, limit)
	if entity.Pagination.Template != "" {
		rendered, err := applyPaginationTemplate(entity, entity.Pagination.Template, entries)
		if err != nil {
			return nil, nil
		}
		return []byte(rendered), nil
	}

	return json.Marshal(entries)
}

func applyPaginationTemplate(entity Entity, template string, entries store.PaginatedEntries) (string, error) {
	jsonEntries, err := json.Marshal(entries.Entries)
	if err != nil {
		return "", err
	}
	attributes := map[string]string{
		"entries":            string(jsonEntries),
		"entity.name":        "entity.Name",
		"pagination.limit":   fmt.Sprint(entries.Limit),
		"pagination.total":   fmt.Sprint(entries.Total),
		"pagination.pages":   fmt.Sprint(entries.Pages),
		"pagination.current": fmt.Sprint(entries.Current),
		"pagination.prev":    fmt.Sprint(entries.Prev),
		"pagination.next":    fmt.Sprint(entries.Next),
	}

	output := template

	for key, value := range attributes {
		output = strings.Replace(output, "${"+key+"}", value, -1)
	}

	return output, nil
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
