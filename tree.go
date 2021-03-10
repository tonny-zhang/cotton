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
		root      *node
		nodeDepth map[int][]*node
	}
	node struct {
		key        string
		nodeType   int
		paramName  string
		fullpath   string
		children   map[string]*node
		isRealNode bool
		handler    HandlerFunc
		middleware []HandlerFunc
	}
	resultFind struct {
		node   *node
		params map[string]string
	}
)

func (n *node) addChild(key string) (*node, error) {
	child, ok := n.children[key]
	if !ok {
		child = newNode(key)
		// for _, c := range n.children {
		// 	c.children
		// }
		n.children[key] = child
	}
	return child, nil
}
func (t *tree) Add(path string, handler HandlerFunc) *node {
	if len(path) == 0 || path[0] != '/' {
		panic(fmt.Errorf("path [%s] must start with /", path))
	}
	var nodeCurrent = t.root
	var start, lenstr = 1, len(path)
	var depth = 1
	if path != nodeCurrent.key {
		checkResult := t.find(path, true)
		if nil != checkResult {
			panic(fmt.Errorf("[%s] conflicts with [%s]", path, checkResult.node.fullpath))
		}
	Search:
		for {
			for i := start; i < lenstr; i++ {
				if path[i] == '/' {
					key := path[start:i]
					if path[start] == '*' {
						panic(fmt.Errorf("action [%s] must end of path [%s]", key, path))
					}
					n, ok := nodeCurrent.children[key]
					if !ok {
						n = newNode(key)
						nodeCurrent.children[key] = n
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
			if path[start] == '*' && len(nodeCurrent.children) > 0 {
				cChild := 0
				for _, c := range nodeCurrent.children {
					if c.key != "/" {
						cChild++
					}
				}
				if cChild > 0 {
					panic(fmt.Errorf("path [%s] conflicts with other rule", path))
				}
			}
			n, ok := nodeCurrent.children[key]
			if !ok {
				n = newNode(key)
				nodeCurrent.children[key] = n
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

func (n *node) match(key string) (isMatch bool) {
	if len(n.key) > 0 && (n.key[0] == ':' || n.key[0] == '*') {
		isMatch = true
	} else if n.key == key {
		isMatch = true
	}
	return
}
func (t *tree) Find(path string) *resultFind {
	return t.find(path, false)
}
func (t *tree) find(path string, isFromAdd bool) *resultFind {
	if path == "/" {
		if t.root.isRealNode {
			return &resultFind{
				node: t.root,
			}
		}
		return nil
	}
	var start, lenstr = 1, len(path)
	var res []string
Split:
	for {
		for i := start; i < lenstr; i++ {
			if path[i] == '/' {
				res = append(res, path[start:i])
				start = i + 1
				continue Split
			}
		}

		if start == lenstr {
			start--
		}
		res = append(res, path[start:])
		break
	}

	var queue = make([]resultFind, 0)
	for _, n := range t.root.children {
		queue = append(queue, resultFind{
			node:   n,
			params: make(map[string]string),
		})
	}
	lastIndex := len(res) - 1
	for i, key := range res {
		var queueTemp = make([]resultFind, 0)
		for _, rf := range queue {
			n := rf.node
			isMatch := n.match(key)
			if !isMatch && len(key) > 0 && key[0] == ':' && n.key != "/" {
				isMatch = true
			}
			if isMatch {
				if n.isRealNode {
					if n.nodeType == nodeCatchAll {
						rf.params[n.paramName] = strings.Join(res[i:], "/")
						return &rf
					}

					if i == lastIndex {
						if n.nodeType == nodeParam {
							rf.params[n.paramName] = key
						}
						return &rf
					}
				}
				var params = make(map[string]string)
				for k, v := range rf.params {
					params[k] = v
				}
				if n.nodeType == nodeParam {
					params[n.paramName] = key
				}
				for _, cn := range n.children {
					queueTemp = append(queueTemp, resultFind{
						node:   cn,
						params: params,
					})
				}
			}
		}
		queue = queueTemp
	}

	return nil
}
func newTree() *tree {
	t := new(tree)
	t.root = newNode("/")
	// t.nodeDepth = make(map[int][]*node)
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

	// fmt.Printf("add node [%s]\n", key)
	return n
}
