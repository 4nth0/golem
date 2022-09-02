package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"
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
	for _, option := range options {
		option(db)
	}

	db.length = 0
	db.Entries = []interface{}{}

	return db
}

func WithLocalFileSync(path string) func(*Database) {
	return func(db *Database) {
		db.FilePath = path
		db.sync = true
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
	output := PaginatedEntries{
		Entries: db.Entries[page*limit : (page+1)*limit],
		Total:   db.length,
		Pages:   int(math.Ceil(float64(db.length) / float64(limit))),
		Current: page,
		Prev:    page - 1,
		Next:    page + 1,
		Limit:   limit,
	}

	if page == 0 {
		output.Prev = 0
	}
	if page == output.Pages {
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

	b, _ := json.MarshalIndent(db.Entries, "", "  ")
	err := ioutil.WriteFile(db.FilePath, b, 0644)
	if err != nil {
		fmt.Println("Err: ", err)
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
