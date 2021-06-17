package cotton

import (
	"fmt"
	"strings"
	"sync"
)

const (
	nodeStatic = iota
	nodeEmpty
	nodeParam
	nodeCatchAll

	keyParam = ":_"
)

var paramsPool sync.Pool
var maxNumParams = 0

type (
	tree struct {
		root *node
	}
	node struct {
		key      string
		fullpath string
		nodeType byte
		children map[string]*node

		isRealNode bool
		handler    HandlerFunc
		middleware []HandlerFunc

		paramName string
	}
	resultFind struct {
		node   *node
		params map[string]string
	}
)

func init() {
	paramsPool.New = func() interface{} {
		return make(map[string]string)
	}
}
func (t *tree) add(path string, handler HandlerFunc) *node {
	if len(path) == 0 || path[0] != '/' {
		panic(fmt.Errorf("path [%s] must start with /", path))
	}
	for i, j := 0, len(path); i < j; i++ {
		if path[i] == '*' {
			pathSub := path[i+1:]
			for ii, jj := 0, len(pathSub); ii < jj; ii++ {
				if pathSub[ii] == '/' {
					panic(fmt.Errorf("action [*%s] must end of path [%s]", pathSub[:ii], path))
				}
			}
			break
		}
	}

	var nodeCurrent = t.root
	var numParams = 0
	var start, lenstr = 1, len(path)
	if path != nodeCurrent.key {
		for i := start; i < lenstr; i++ {
			if path[i] == '/' {
				key := path[start:i]

				nodeCurrent.insertNode(key, path)
				child, ok := nodeCurrent.children[key]
				if !ok {
					child = nodeCurrent.children[keyParam]
					numParams++
				} else {
					for _, nc := range nodeCurrent.children {
						if nc.nodeType == nodeCatchAll {
							panic(fmt.Errorf("path [%s] conflicts with [%s]", path, nc.fullpath))
						}
					}
				}
				nodeCurrent = child
				start = i + 1
				continue
			}
		}

		if start == lenstr {
			start--
		}
		key := path[start:]
		nodeCurrent.insertNode(key, path)
		child, ok := nodeCurrent.children[key]
		if !ok {
			child = nodeCurrent.children[keyParam]

			for _, c := range nodeCurrent.children {
				if c.isRealNode && c.nodeType == nodeStatic && c.key != "/" {
					panic(fmt.Errorf("path [%s] conflicts with [%s]", path, c.fullpath))
				}
			}
			numParams++
		} else {
			if c, ok := nodeCurrent.children[keyParam]; ok && c.isRealNode {
				panic(fmt.Errorf("path [%s] conflicts with [%s]", path, c.fullpath))
			}
		}

		nodeCurrent = child

	}
	nodeCurrent.isRealNode = true
	nodeCurrent.handler = handler
	nodeCurrent.fullpath = path

	if maxNumParams < numParams {
		maxNumParams = numParams

		paramsPool.New = func() interface{} {
			return make(map[string]string, maxNumParams)
		}
	}
	return nodeCurrent
}

func (n *node) print(deep int) {
	key := n.key
	if n.nodeType == nodeParam {
		key = ":" + n.paramName
	} else if n.nodeType == nodeCatchAll {
		key = "*" + n.paramName
	}
	fmt.Printf("%s %d %s %v\n", strings.Repeat("    ", deep), deep, key, n.isRealNode)
	for _, n := range n.children {
		n.print(deep + 1)
	}
}

func (n *node) insertNode(key string, fullpath string) {
	if n.nodeType == nodeCatchAll {
		panic(fmt.Errorf("action [%s] must end of path [%s]", n.key, fullpath))
	}
	keyCheck := key
	if keyCheck[0] == ':' || keyCheck[0] == '*' {
		keyCheck = keyParam

		if cCheck, ok := n.children[keyParam]; ok && cCheck.paramName != key[1:] {
			var fullpathV = cCheck.fullpath
			if !cCheck.isRealNode {
				for _, v := range n.children {
					if v.key == "/" {
						continue
					}

					vCheck := v
					for !vCheck.isRealNode {
						for _, vv := range v.children {
							vCheck = vv
							break
						}
					}
					fullpathV = vCheck.fullpath

				}
			}

			panic(fmt.Errorf("[%s] in path [%s] conflicts with [%s]", key, fullpath, fullpathV))
		}
	}

	if _, ok := n.children[keyCheck]; !ok {
		n.children[keyCheck] = newNode(key)
	}
}
func (n *node) find(path string) (result resultFind) {
	if path == "/" && n.key == path {
		if n.isRealNode {
			result.node = n
		}
		return
	}

	var child *node
	var ok bool
	var start, lenstr = 1, len(path)

	result.params = paramsPool.Get().(map[string]string)

	var keyPrev = ""
	for i := start; i < lenstr; i++ {
		if path[i] == '/' {
			key := path[start:i]
			if key == keyPrev {
				start = i + 1
				continue
			}
			keyPrev = key

			child, ok = n.children[key]
			if !ok {
				child, ok = n.children[keyParam]
				if ok {
					if child.nodeType == nodeCatchAll {
						result.params[child.paramName] = path[start:]
						paramsPool.Put(result.params)
						result.node = child
						return
					}
					result.params[child.paramName] = key
				}
			}
			if !ok {
				paramsPool.Put(result.params)
				return
			}
			n = child
			start = i + 1
		}
	}
	if start == lenstr {
		start--
	}
	key := path[start:]
	child, ok = n.children[key]
	if !ok {
		child, ok = n.children[keyParam]
		if ok {
			result.params[child.paramName] = key
		}
	}
	if ok && child.isRealNode {
		result.node = child
	}
	paramsPool.Put(result.params)
	return
}

func newTree() *tree {
	t := new(tree)
	n := newNode("/")
	t.root = n
	return t
}
func newNode(key string) *node {
	n := new(node)
	n.nodeType = nodeEmpty
	if len(key) > 0 {
		switch key[0] {
		case ':':
			n.nodeType = nodeParam
			n.paramName = key[1:]
			n.key = keyParam
		case '*':
			n.nodeType = nodeCatchAll
			n.paramName = key[1:]
			n.key = keyParam
		default:
			n.nodeType = nodeStatic
			n.key = key
		}

	}
	n.children = make(map[string]*node)
	return n
}
