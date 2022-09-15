package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Entries []interface{}

type Database struct {
	sync     bool          `json:"-"`
	FilePath string        `json:"-"`
	length   int           `json:"-"`
	Entries  []interface{} `json:"entries"`
	mux      sync.Mutex    `json:"-"`
}

type Option struct{}

func New(options ...func(*Database)) *Database {
	db := &Database{}
	db.length = 0
	db.Entries = make([]interface{}, 0)

	for _, option := range options {
		option(db)
	}

	return db
}

func WithLocalFile(path string, sync bool) func(*Database) {
	return func(db *Database) {
		db.FilePath = path
		db.sync = sync
	}
}

func WithData(data Entries) func(*Database) {
	return func(db *Database) {
		db.length = len(data)
		db.Entries = data
	}
}

func (db *Database) List() Entries {
	return db.Entries
}

type PaginatedEntries struct {
	Entries Entries `json:"entries"`
	Limit   int     `json:"limit"`
	Total   int     `json:"total"`
	Pages   int     `json:"pages"`
	Current int     `json:"current"`
	Prev    int     `json:"prev"`
	Next    int     `json:"next"`
}

func (db *Database) PaginatedList(page, limit int) PaginatedEntries {
	if limit > db.length {
		limit = db.length
	}
	if db.length == 0 || page > db.length/limit {
		return PaginatedEntries{}
	}
	if page < 0 {
		page = 0
	}
	pages := int(math.Ceil(float64(db.length) / float64(limit)))
	output := PaginatedEntries{
		Entries: db.Entries[page*limit : (page+1)*limit],
		Total:   db.length,
		Pages:   pages,
		Current: page,
		Prev:    page - 1,
		Next:    page + 1,
		Limit:   limit,
	}

	if page == 0 {
		output.Prev = 0
	}
	if page == pages {
		output.Next = page
	}

	return output
}

func (db *Database) Load() error {
	_, err := os.Stat(db.FilePath)

	if err == nil {
		data, err := ioutil.ReadFile(db.FilePath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &db.Entries)
		db.length = len(db.Entries)
		return err
	}

	return err
}

func (db *Database) Push(entry interface{}) error {
	db.mux.Lock()
	db.Entries = append(db.Entries, entry)
	db.length++
	db.Save()
	db.mux.Unlock()
	return nil
}

func (db *Database) Save() {
	if !db.sync {
		return
	}

	b, _ := json.Marshal(db.Entries)
	err := ioutil.WriteFile(db.FilePath, b, 0644)
	if err != nil {
		log.Error("Err: ", err)
	}
}

func (db *Database) GetByIndex(index int) (interface{}, error) {
	if index >= db.length || index < 0 {
		return nil, errors.New("NOT_FOUND")
	}
	return db.Entries[index], nil
}

func (db *Database) DeleteFromIndex(index int) error {
	if index >= db.length || index < 0 {
		return errors.New("NOT_FOUND")
	}
	db.mux.Lock()
	db.Entries = append(db.Entries[:index], db.Entries[index+1:]...)
	db.length--
	db.Save()
	db.mux.Unlock()
	return nil
}
