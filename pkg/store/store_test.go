package store

import (
	"fmt"
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

	db.Load()

	fmt.Println(db.List())

	// DB should have a length of 8
	assert.Equal(t, usersLentgh, db.length)

	// DB should have an entries length of 8
	assert.Len(t, db.Entries, usersLentgh)

	// DB value of path should be equal to "./test/path"
	assert.Equal(t, db.FilePath, usersGoldenPath)

	// DB value of sync should be equal to false
	assert.False(t, db.sync)

}
