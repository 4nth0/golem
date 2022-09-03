package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/4nth0/golem/store"
)

func renderPaginatedList(entityName string, entityConfig Entity, r *http.Request, db *store.Database) ([]byte, error) {
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
	if entityConfig.Pagination.Template != "" {
		rendered, err := applyPaginationTemplate(entityName, entityConfig.Pagination.Template, entries)
		if err != nil {
			return nil, nil
		}
		return []byte(rendered), nil
	}

	return json.Marshal(entries)
}

func applyPaginationTemplate(entityName string, template string, entries store.PaginatedEntries) (string, error) {
	jsonEntries, err := json.Marshal(entries.Entries)
	if err != nil {
		return "", err
	}
	attributes := map[string]string{
		"entries":            string(jsonEntries),
		"entity.name":        entityName,
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
