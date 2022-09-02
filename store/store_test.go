package store

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	usersGoldenPath = "../test/users.db.golden.json"
	usersLentgh     = 8
)

func Test_Init(t *testing.T) {
	path := "./test/path"
	db := New(path, false)

	// DB should start with a length of 0
	assert.Equal(t, db.length, 0)

	// DB should start with an entries length of 0
	assert.Len(t, db.Entries, 0)

	// DB value of path should be equal to "./test/path"
	assert.Equal(t, db.FilePath, path)

	// DB value of sync should be equal to false
	assert.False(t, db.sync)
}

func Test_Load(t *testing.T) {
	db := New(usersGoldenPath, false)

	err := db.Load()

	assert.Nil(t, err)
	// DB should have a length of 8
	assert.Equal(t, usersLentgh, db.length)

	// DB should have an entries length of 8
	assert.Len(t, db.Entries, usersLentgh)

	db2 := New("./path/to/nonexistent/file", false)
	err = db2.Load()

	assert.NotNil(t, err)
}

func Test_List(t *testing.T) {
	db := New(usersGoldenPath, false)
	err := db.Load()
	assert.Nil(t, err)

	entries := db.List()

	assert.Len(t, entries, usersLentgh)
}

func Test_PaginatedList(t *testing.T) {
	db := New(usersGoldenPath, false)
	err := db.Load()
	assert.Nil(t, err)

	entries := db.PaginatedList(0, 4)

	assert.Len(t, entries.Entries, 4)
	assert.Equal(t, entries.Limit, 4)
	assert.Equal(t, entries.Total, usersLentgh)
	assert.Equal(t, entries.Pages, usersLentgh/4)
	assert.Equal(t, entries.Current, 0)
	assert.Equal(t, entries.Prev, 0)
	assert.Equal(t, entries.Next, 1)

	entries = db.PaginatedList(-1, 4)

	assert.Len(t, entries.Entries, 4)
	assert.Equal(t, entries.Limit, 4)
	assert.Equal(t, entries.Total, usersLentgh)
	assert.Equal(t, entries.Pages, usersLentgh/4)
	assert.Equal(t, entries.Current, 0)
	assert.Equal(t, entries.Prev, 0)
	assert.Equal(t, entries.Next, 1)
}

func Test_Push(t *testing.T) {
	db := New(usersGoldenPath, false)
	err := db.Load()
	assert.Nil(t, err)

	assert.Equal(t, usersLentgh, db.length)

	err = db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)
	assert.Nil(t, err)

	assert.Equal(t, usersLentgh+1, db.length)
}

func Test_GetByIndex(t *testing.T) {
	db := New(usersGoldenPath, false)

	err := db.Load()
	assert.Nil(t, err)
	entries := db.List()
	entry, _ := db.GetByIndex(3)
	_, err = db.GetByIndex(30)

	assert.Equal(t, entries[3], entry)
	assert.NotNil(t, err)
}

func Test_DeleteFromIndex(t *testing.T) {
	db := New(usersGoldenPath, false)

	err := db.Load()
	assert.Nil(t, err)

	toBeDeleted, _ := db.GetByIndex(3)

	assert.Equal(t, usersLentgh, db.length)

	err = db.DeleteFromIndex(3)
	assert.Nil(t, err)

	entry, _ := db.GetByIndex(3)
	entries := db.List()

	assert.Equal(t, usersLentgh-1, db.length)
	assert.Len(t, entries, usersLentgh-1)
	assert.NotEqual(t, entry, toBeDeleted)

	err = db.DeleteFromIndex(30)
	assert.NotNil(t, err)
}

func Test_Save(t *testing.T) {
	file, err := ioutil.TempFile("", "golem.test-sync.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	db := New(file.Name(), true)

	db.Load() //nolint:all

	err = db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)
	assert.Nil(t, err)

	err = db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)
	assert.Nil(t, err)

	err = db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)
	assert.Nil(t, err)

	err = db.DeleteFromIndex(1)
	assert.Nil(t, err)

	db2 := New(file.Name(), true)
	err = db2.Load()
	assert.Nil(t, err)

	assert.Equal(t, 2, db2.length)
}
