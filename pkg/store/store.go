package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
)

type Entries []interface{}

type Database struct {
	Sync     bool          `json:"-"`
	FilePath string        `json:"-"`
	Length   int           `json:"length"`
	Entries  []interface{} `json:"entries"`
	mux      sync.Mutex    `json:"-"`
}

func (db *Database) InitDefault() {
	db.Length = 0
	db.Entries = []interface{}{}
}

func (db Database) List() Entries {
	return db.Entries
}

func (db *Database) Push(entry interface{}) error {
	db.mux.Lock()
	db.Entries = append(db.Entries, entry)
	db.Length++
	db.Save()
	db.mux.Unlock()
	return nil
}

func (db *Database) Save() {
	if !db.Sync {
		return
	}

	b, _ := json.MarshalIndent(db, "", "  ")
	err := ioutil.WriteFile(db.FilePath, b, 0644)
	if err != nil {
		fmt.Println("Err: ", err)
	}
}

func (db Database) GetByIndex(index int) (interface{}, error) {
	if index >= db.Length || index < 0 {
		return nil, errors.New("NOT_FOUND")
	}
	return db.Entries[index], nil
}

func (db *Database) DeleteFromIndex(index int) error {
	if index >= db.Length || index < 0 {
		return errors.New("NOT_FOUND")
	}
	db.mux.Lock()
	db.Entries = append(db.Entries[:index], db.Entries[index+1:]...)
	db.Length--
	db.Save()
	db.mux.Unlock()
	return nil
}
