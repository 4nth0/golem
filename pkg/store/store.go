package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func New(path string, sync bool) *Database {
	db := Database{
		FilePath: path,
		sync:     sync,
	}

	db.length = 0
	db.Entries = []interface{}{}

	fmt.Println("Init")

	return &db
}

func (db Database) List() Entries {
	return db.Entries
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

func (db Database) GetByIndex(index int) (interface{}, error) {
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
