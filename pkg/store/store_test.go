package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	usersGoldenPath = "../../test/users.db.golden.json"
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
	db.Load()

	entries := db.List()

	assert.Len(t, entries, usersLentgh)
}

func Test_Push(t *testing.T) {
	db := New(usersGoldenPath, false)
	db.Load()

	assert.Equal(t, usersLentgh, db.length)

	db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)

	assert.Equal(t, usersLentgh+1, db.length)
}

func Test_GetByIndex(t *testing.T) {
	db := New(usersGoldenPath, false)

	db.Load()
	entries := db.List()
	entry, _ := db.GetByIndex(3)
	_, err := db.GetByIndex(30)

	assert.Equal(t, entries[3], entry)
	assert.NotNil(t, err)
}

func Test_DeleteFromIndex(t *testing.T) {
	db := New(usersGoldenPath, false)

	db.Load()

	toBeDeleted, _ := db.GetByIndex(3)

	assert.Equal(t, usersLentgh, db.length)

	db.DeleteFromIndex(3)

	entry, _ := db.GetByIndex(3)
	entries := db.List()

	assert.Equal(t, usersLentgh-1, db.length)
	assert.Len(t, entries, usersLentgh-1)
	assert.NotEqual(t, entry, toBeDeleted)

	err := db.DeleteFromIndex(30)
	assert.NotNil(t, err)
}

func Test_Save(t *testing.T) {
	path := usersGoldenPath + ".test-sync"
	db := New(path, true)

	db.Load()

	db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)
	db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)
	db.Push(`{"name": "Jody Mills", "type": "Hunter"}`)

	db.DeleteFromIndex(1)

	db2 := New(path, true)
	db2.Load()

	assert.Equal(t, 2, db2.length)

	os.Remove(path)
}
