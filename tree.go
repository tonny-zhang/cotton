package cotton

import (
	"fmt"
	"strconv"
)

type (
	tree struct {
		root *node
	}
	node struct {
		key        string
		fullpath   string
		children   map[string]*node
		isRealNode bool
		handler    HandlerFunc
	}
	resultFind struct {
		node   *node
		params map[string]string
	}
)

func run(n *node, dep int) {
	if dep == 1 {
		fmt.Printf("%d key = [%s] fullpath = %s isHandle = %v\n", dep, n.key, "", n.isRealNode)
	} else {
		fmt.Printf("%"+(strconv.Itoa(dep*5))+"s %d key = [%s] fullpath = %s isHandle = %v\n", "     ", dep, n.key, "", n.isRealNode)
	}

	for _, c := range n.children {
		run(c, dep+1)
	}
}
func (t *tree) print() {
	fmt.Println("==================")

	run(t.root, 1)
	fmt.Println("==================")
}
func (t *tree) Add(path string, handler HandlerFunc) {
	var nodeCurrent = t.root
	var start, lenstr = 0, len(path)
	if path != "/" {
	Search:
		for {
			for i := start; i < lenstr; i++ {
				if path[i] == '/' {
					key := path[start:i]
					n, ok := nodeCurrent.children[key]
					if !ok {
						n = newNode(key)
						nodeCurrent.children[key] = n
					}

					nodeCurrent = n
					start = i + 1
					continue Search
				}
			}

			key := path[start:]
			n, ok := nodeCurrent.children[key]
			if !ok {
				n = newNode(key)
				nodeCurrent.children[key] = n
			}
			nodeCurrent = n
			break
		}
	}
	nodeCurrent.isRealNode = true
	nodeCurrent.handler = handler
}

func (n *node) match(key string) (isMatch bool, paramName string) {
	if len(n.key) > 0 && n.key[0] == ':' {
		isMatch = true
		paramName = n.key[1:]
	} else if n.key == key {
		isMatch = true
	}
	return
}
func (t *tree) Find(path string) *resultFind {
	if path == "/" {
		if t.root.isRealNode {
			return &resultFind{
				node: t.root,
			}
		}
		return nil
	}
	var start, lenstr = 0, len(path)

	var res []string
	{
	Split:
		for {
			for i := start; i < lenstr; i++ {
				if path[i] == '/' {
					res = append(res, path[start:i])
					start = i + 1
					continue Split
				}
			}

			res = append(res, path[start:])
			break
		}
	}

	res = res[1:]
	// var queue = make([]resultFind, 0)
	// for _, n := range t.root.children[""].children {
	// 	queue = append(queue, resultFind{
	// 		node:   n,
	// 		params: make(map[string]string),
	// 	})
	// }
	// for i, key := range res {
	// 	var queueTemp = make([]resultFind, 0)
	// 	for _, rf := range queue {
	// 		n := rf.node
	// 		isMatch, name := n.match(key)
	// 		if isMatch {
	// 			if i == len(res)-1 && n.isRealNode {
	// 				if name != "" {
	// 					rf.params[name] = key
	// 				}
	// 				return &rf
	// 			}
	// 			var params = make(map[string]string)
	// 			for k, v := range rf.params {
	// 				params[k] = v
	// 			}
	// 			if name != "" {
	// 				params[name] = key
	// 			}
	// 			for _, cn := range n.children {
	// 				queueTemp = append(queueTemp, resultFind{
	// 					node:   cn,
	// 					params: params,
	// 				})
	// 			}
	// 		}
	// 	}
	// 	queue = queueTemp
	// }

	return nil
}
func newTree() *tree {
	t := new(tree)
	t.root = newNode("/")
	return t
}
func newNode(key string) *node {
	n := new(node)
	n.key = key
	n.children = make(map[string]*node)

	// fmt.Printf("add node [%s]\n", key)
	return n
}
