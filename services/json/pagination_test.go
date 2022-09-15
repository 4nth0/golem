package json

import (
	"testing"

	"github.com/4nth0/golem/store"
	"github.com/stretchr/testify/assert"
)

const (
	template = `{
		"entries": ${entries},
		"_metadata": {
		  "per_page": ${pagination.limit},
		  "page": ${pagination.current},
		  "page_count": ${pagination.pages},
		  "total_count": ${pagination.total},
		  "links": {
			"self":  "/${entity.name}?limit=${pagination.limit}&page=${pagination.current}",
			"first": "/${entity.name}?limit=${pagination.limit}&page=${pagination.current}",
			"next":  "/${entity.name}?limit=${pagination.limit}&page=${pagination.next}",
			"prev":  "/${entity.name}?limit=${pagination.limit}&page=${pagination.prev}"
		  }
		}
	  }`
	expectedResult = `{
		"entries": ["lorem","ipsum","dolor","sit","amet","consectetur"],
		"_metadata": {
		  "per_page": 2,
		  "page": 0,
		  "page_count": 3,
		  "total_count": 6,
		  "links": {
			"self":  "/users?limit=2&page=0",
			"first": "/users?limit=2&page=0",
			"next":  "/users?limit=2&page=1",
			"prev":  "/users?limit=2&page=0"
		  }
		}
	  }`
)

func Test_applyPaginationTemplate(t *testing.T) {
	entityName := "users"
	entries := store.Entries{
		"lorem",
		"ipsum",
		"dolor",
		"sit",
		"amet",
		"consectetur",
	}
	paginatedEntries := store.PaginatedEntries{
		Entries: entries,
		Limit:   2,
		Total:   6,
		Pages:   3,
		Current: 0,
		Prev:    0,
		Next:    1,
	}

	rendered, err := applyPaginationTemplate(entityName, template, paginatedEntries)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, rendered)
}
