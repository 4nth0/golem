package tree

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func DryHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {}

func TestAddNode(t *testing.T) {
	tree := NewTree()

	tree.AddNode("/path/to/heaven", "GET", DryHandler)
	node := tree.Childs["path"].Childs["to"].Childs["heaven"]
	h := node.Handler["GET"]

	assert.NotNil(t, h)
	assert.Len(t, node.Childs, 0)
	assert.Equal(t, "", node.VarName)

	tree.AddNode("/path/to/:param", "GET", DryHandler)
	node = tree.Childs["path"].Childs["to"].Childs["*"]
	h = node.Handler["GET"]

	assert.NotNil(t, h)
	assert.Len(t, node.Childs, 0)
	assert.Equal(t, "param", node.VarName)

	assert.Len(t, tree.Childs, 1)
	assert.Len(t, tree.Childs["path"].Childs["to"].Childs, 2)
}

func TestGetNode(t *testing.T) {
	tree := NewTree()

	tree.AddNode("/path/to/heaven", "GET", DryHandler)
	tree.AddNode("/path/to/:param", "GET", DryHandler)

	handler, params, _ := tree.GetNode("/path/to/heaven", "GET")
	assert.NotNil(t, handler)
	assert.Equal(t, params, map[string]string{})

	handler, params, _ = tree.GetNode("/path/to/hell", "GET")
	assert.NotNil(t, handler)
	assert.Equal(t, "hell", params["param"])

	handler, params, _ = tree.GetNode("/unexistant/path", "GET")
	assert.Nil(t, handler)
	assert.Nil(t, params)

	_, _, err := tree.GetNode("/path/to/heaven", "POST")
	assert.NotNil(t, err)

}

func TestNodeWithDifferentMethods(t *testing.T) {
	tree := NewTree()

	tree.AddNode("/path/to/resource", "GET", DryHandler)
	tree.AddNode("/path/to/resource", "POST", DryHandler)

	// Vérifie la présence des deux méthodes
	getHandler, _, _ := tree.GetNode("/path/to/resource", "GET")
	assert.NotNil(t, getHandler)

	postHandler, _, _ := tree.GetNode("/path/to/resource", "POST")
	assert.NotNil(t, postHandler)
}

func TestTrailingSlash(t *testing.T) {
	tree := NewTree()

	tree.AddNode("/path/with/slash/", "GET", DryHandler)

	// Vérifie la compatibilité avec et sans le slash final
	handlerWithSlash, _, _ := tree.GetNode("/path/with/slash/", "GET")
	assert.NotNil(t, handlerWithSlash)

	handlerWithoutSlash, _, _ := tree.GetNode("/path/with/slash", "GET")
	assert.NotNil(t, handlerWithoutSlash)
}
