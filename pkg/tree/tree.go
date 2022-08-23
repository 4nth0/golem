package tree

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
)

const PathDelimiter = "/"
const VarPrefix = ":"
const Wildcard = "*"

type TreeNode struct {
	Handler map[string]Handler
	Childs  map[string]*TreeNode
	VarName string
	Parent  *TreeNode  `json:"-"`
	mux     sync.Mutex `json:"-"`
}

type Handler func(w http.ResponseWriter, r *http.Request, params map[string]string)

func NewTree() *TreeNode {
	return &TreeNode{
		Handler: nil,
		Childs:  map[string]*TreeNode{},
	}
}

func (t *TreeNode) AddNode(path string, method string, handler Handler) {
	path = strings.TrimPrefix(path, PathDelimiter)
	splitted := strings.Split(path, PathDelimiter)

	currentNode := t

	t.mux.Lock()
	defer t.mux.Unlock()

	for i := 0; i < len(splitted); i++ {
		key := splitted[i]
		varName := ""

		if strings.HasPrefix(key, VarPrefix) {
			varName = key[len(VarPrefix):]
			key = Wildcard
		}

		if _, ok := currentNode.Childs[key]; !ok {
			currentNode.Childs[key] = &TreeNode{
				Handler: nil,
				Childs:  map[string]*TreeNode{},
				Parent:  currentNode,
			}
		}

		currentNode = currentNode.Childs[key]
		if varName != "" {
			currentNode.VarName = varName
		}

		if i == len(splitted)-1 {
			if currentNode.Handler == nil {
				currentNode.Handler = map[string]Handler{}
			}
			currentNode.Handler[method] = handler
			return
		}
	}
}

func (t *TreeNode) GetNode(path, method string) (Handler, map[string]string, error) {
	path = strings.TrimPrefix(path, PathDelimiter)
	splitted := strings.Split(path, PathDelimiter)
	params := map[string]string{}
	currentNode := t

	for i := 0; i < len(splitted); i++ {
		key := splitted[i]

		if _, ok := currentNode.Childs[key]; !ok {
			if _, ok := currentNode.Childs[Wildcard]; !ok {
				return nil, nil, nil
			} else {
				params[currentNode.Childs[Wildcard].VarName] = key
				key = Wildcard
			}
		}

		if i == len(splitted)-1 {
			if currentNode.Childs[key].Handler == nil {
				return nil, nil, nil
			} else {
				if _, ok := currentNode.Childs[key].Handler[method]; !ok {
					return nil, nil, errors.New("METHOD_NOT_ALLOWED")
				}
				return currentNode.Childs[key].Handler[method], params, nil
			}
		}

		currentNode = currentNode.Childs[key]
	}

	return nil, nil, nil
}

func (t *TreeNode) RemoveNode(path string) {
	path = strings.TrimPrefix(path, PathDelimiter)
	splitted := strings.Split(path, PathDelimiter)

	currentNode := t

	for i := 0; i < len(splitted); i++ {
		key := splitted[i]

		if _, ok := currentNode.Childs[key]; !ok {
			return
		}

		if i == len(splitted)-1 {
			if len(currentNode.Childs[key].Childs) == 0 {
				delete(currentNode.Childs, key)
			} else {
				currentNode.Childs[key].Handler = nil
			}
			return
		}

		currentNode = currentNode.Childs[key]
	}
}
func (t *TreeNode) Dump() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}
