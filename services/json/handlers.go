package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/4nth0/golem/log"
	"github.com/4nth0/golem/router/tree"
	"github.com/4nth0/golem/store"
)

const (
	ContentTypeJSON         = "application/json"
	ParamIndexKey           = "index"
	DefaultPaginationtLimit = 10
	DefaultPaginationtPage  = 0
)

var empty []byte = []byte("{}")

func GetHandler(db *store.Database) tree.Handler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", ContentTypeJSON)

		index, _ := strconv.Atoi(params[ParamIndexKey])
		entry, err := db.GetByIndex(index)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write(empty)
			log.Error("Unable to write response.", "err", err)
			return
		}

		list, err := json.Marshal(entry)
		if err != nil {
			fmt.Println("Err: ", err)
		}

		_, err = w.Write(list)
		if err != nil {
			log.Error("Unable to write response.", "err", err)
		}
	}
}

func ListHandler(db *store.Database, entityName string, entity Entity) tree.Handler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", ContentTypeJSON)

		var list []byte
		var err error

		if entity.Pagination != nil {
			list, err = renderPaginatedList(entityName, entity, r, db)
			if err != nil {
				log.Error("Err: ", err)
			}
		} else {
			list, err = json.Marshal(db.List())
			if err != nil {
				log.Error("Err: ", err)
			}
		}

		_, err = w.Write(list)
		if err != nil {
			log.Error("Unable to write response.", "err", err)
		}
	}
}

func CreateHandler(db *store.Database) tree.Handler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		decoder := json.NewDecoder(r.Body)
		var t interface{}
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		err = db.Push(t)
		if err != nil {
			log.Error("Unable to push request in db", "err", err)
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteHandler(db *store.Database) tree.Handler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		index, _ := strconv.Atoi(params["index"])
		err := db.DeleteFromIndex(index)
		if err != nil {
			log.Error("Unable to delete entry", "err", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
