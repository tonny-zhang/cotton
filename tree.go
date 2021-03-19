package cotton

import (
	"fmt"
	"strings"
)

const (
	nodeStatic = iota
	nodeEmpty
	nodeParam
	nodeCatchAll
)

type (
	tree struct {
		root *node
	}
	node struct {
		key      string
		fullpath string
		nodeType int
		children map[string]*node

		isRealNode bool
		handler    HandlerFunc
		middleware []HandlerFunc

		paramName string
		paramKey  string
	}
	resultFind struct {
		node   *node
		params map[string]string
	}
)

func (n *node) insertNode(key string, fullpath string) *node {
	child := newNode(key)
	if n.nodeType == nodeCatchAll {
		panic(fmt.Errorf("action [%s] must end of path [%s]", n.key, fullpath))
	}
	if child.nodeType == nodeParam || child.nodeType == nodeCatchAll {
		for _, v := range n.children {
			if v.key == "/" || v.key == "a" {
				continue
			}

			vCheck := v
			for !vCheck.isRealNode {
				for _, vv := range v.children {
					vCheck = vv
					break
				}
			}
			fullpathV := vCheck.fullpath
			panic(fmt.Errorf("[%s] in path [%s] conflicts with [%s]", key, fullpath, fullpathV))
		}
		n.paramKey = key
	}
	n.children[key] = child

	return child
}
func (t *tree) Add(path string, handler HandlerFunc) *node {
	if len(path) == 0 || path[0] != '/' {
		panic(fmt.Errorf("path [%s] must start with /", path))
	}
	var nodeCurrent = t.root
	var start, lenstr = 1, len(path)
	var depth = 1
	if path != nodeCurrent.key {
	Search:
		for {
			for i := start; i < lenstr; i++ {
				if path[i] == '/' {

					key := path[start:i]

					n, ok := nodeCurrent.children[key]
					if !ok {
						n = nodeCurrent.insertNode(key, path)
					}

					nodeCurrent = n
					start = i + 1
					depth++
					continue Search
				}
			}

			if start == lenstr {
				start--
			}
			key := path[start:]
			n, ok := nodeCurrent.children[key]
			if !ok {
				n = nodeCurrent.insertNode(key, path)
			}
			nodeCurrent = n
			depth++
			break
		}
	}
	nodeCurrent.isRealNode = true
	nodeCurrent.handler = handler
	nodeCurrent.fullpath = path
	if handler != nil {
		nodeCurrent.middleware = append(nodeCurrent.middleware, handler)
	}

	return nodeCurrent
}

func (node *node) print(deep int) {
	fmt.Printf("%s %d %s %v %d %s\n", strings.Repeat("    ", deep), deep, node.key, node.isRealNode, node.nodeType, node.paramKey)
	for _, n := range node.children {
		n.print(deep + 1)
	}
}
func (n *node) find(path string) *resultFind {
	if path == "/" && n.key == path {
		if n.isRealNode {
			return &resultFind{
				node: n,
			}
		}
		return nil
	}
	var params = make(map[string]string)
	var start, lenstr = 1, len(path)
	for i := start; i < lenstr; i++ {
		if path[i] == '/' {
			key := path[start:i]

			nc, ok := n.children[key]
			if !ok {
				nc, ok = n.children[n.paramKey]
				if ok {
					if nc.nodeType == nodeCatchAll {
						params[nc.paramName] = path[start:]
						return &resultFind{
							node:   nc,
							params: params,
						}
					}
					params[nc.paramName] = key
				}
			}
			if !ok {
				return nil
			}
			n = nc
			start = i + 1
		}
	}
	if start == lenstr {
		start--
	}
	key := path[start:]
	nc, ok := n.children[key]
	if !ok {
		nc, ok = n.children[n.paramKey]
		if ok {
			params[nc.paramName] = key
		}
	}
	if ok {
		return &resultFind{
			node:   nc,
			params: params,
		}
	}
	return nil
}

func (t *tree) Find(path string) *resultFind {
	return t.root.find(path)
}

func newTree() *tree {
	t := new(tree)
	t.root = newNode("/")
	return t
}
func newNode(key string) *node {
	n := new(node)
	n.key = key
	n.nodeType = nodeEmpty
	if len(key) > 0 {
		switch key[0] {
		case ':':
			n.nodeType = nodeParam
			n.paramName = key[1:]
		case '*':
			n.nodeType = nodeCatchAll
			n.paramName = key[1:]
		default:
			n.nodeType = nodeStatic
		}

	}
	n.children = make(map[string]*node)
	return n
}
