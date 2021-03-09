package tree

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tonny-zhang/cotton/utils"
)

type HandlerFunc func()

var regParamDetail = regexp.MustCompile(":([a-z]+)(<num>|{([^}]+)})?")

type (
	tree struct {
		root     *node
		allNodes map[int][]*node
	}
	resultFind struct {
		params map[string]string
		node   *node
	}
	node struct {
		parent   *node
		children map[string]*node

		key         string
		depth       int
		fullpath    string
		isRealNode  bool
		pattern     *nodePattern
		fullPattern *nodeFullPattern
		handler     HandlerFunc
	}
	nodeFullPattern struct {
		keys []string
		exp  *regexp.Regexp
	}
	nodePattern struct {
		name   string
		expStr string
		exp    *regexp.Regexp
	}
)

func newTree() *tree {
	nRoot := newNode("/", nil)
	nRoot.fullpath = "/"
	allnodes := make(map[int][]*node)
	return &tree{
		root:     nRoot,
		allNodes: allnodes,
	}
}

func newNode(key string, parent *node) *node {
	n := node{
		key:      key,
		parent:   parent,
		children: make(map[string]*node),
	}
	if nil != parent {
		n.depth = parent.depth + 1
	} else {
		n.depth = 1
	}
	if len(key) > 0 && key[0] == ':' {
		resultDetail := regParamDetail.FindAllStringSubmatch(key, -1)[0]
		regexpStr := resultDetail[3]
		isNum := false
		isRegexp := "" != regexpStr

		if !isRegexp {
			isNum = resultDetail[2] == "<num>"
			if isNum {
				regexpStr = "\\d+"
			} else {
				regexpStr = "[^/]+"
			}
		}
		n.pattern = &nodePattern{
			name:   resultDetail[1],
			expStr: "(" + regexpStr + ")",
			exp:    regexp.MustCompile("^" + regexpStr + "$"),
		}
	}
	return &n
}
func (n *node) isLast() bool {
	return len(n.children) == 0
}
func run(n *node, dep int) {
	if dep == 1 {
		fmt.Printf("%d key = %s fullpath = %s isHandle = %v\n", n.depth, n.key, n.fullpath, n.isRealNode)
	} else {
		fmt.Printf("%"+(strconv.Itoa(dep*5))+"s %d key = %s fullpath = %s isHandle = %v\n", "     ", n.depth, n.key, n.fullpath, n.isRealNode)
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
func (t *tree) Add(pattern string, handler HandlerFunc) {
	var currentNode = t.root
	pattern = utils.CleanPath(pattern)
	fullpath := pattern
	if fullpath != currentNode.key {
		pattern = strings.TrimLeft(pattern, "/")
		var expStr string
		var expNames []string
		for _, key := range strings.Split(pattern, "/") {
			n, ok := currentNode.children[key]
			if !ok {
				n = newNode(key, currentNode)
				currentNode.children[key] = n

			}
			if n.pattern == nil {
				expStr += "/" + key
			} else {
				expStr += "/" + n.pattern.expStr
				expNames = append(expNames, n.pattern.name)
			}
			currentNode = n
		}

		currentNode.handler = handler
		currentNode.fullpath = fullpath
		currentNode.isRealNode = true
		currentNode.fullPattern = &nodeFullPattern{
			exp:  regexp.MustCompile("^" + expStr + "$"),
			keys: expNames,
		}
		fmt.Println("exp", expStr, expNames)
		t.allNodes[currentNode.depth] = append(t.allNodes[currentNode.depth], currentNode)
	}

	fmt.Println("add", fullpath, currentNode.isRealNode, currentNode.depth)
}

// /	=> 	["", ""]
// /a	=>	["", "a"]
// /a/	=>	["", "a", ""]
func (n *node) find(path []string) (nodes []*node) {
	lenPath := len(path)
	if lenPath >= 1 {
		// if path[0] == "" {
		// 	path[0] = "/"
		// }
		key := path[0]

		if key == "" || (n.pattern == nil && n.key == key) || (n.pattern != nil) {
			if len(path) == 1 && n.isRealNode {
				nodes = append(nodes, n)
			} else {
				for _, cn := range n.children {
					nodes = append(nodes, cn.find(path[1:])...)
				}
			}
		}
	}
	return
}
func (t *tree) Find(path string) *resultFind {
	// path = utils.CleanPath(path)
	// // fmt.Println("find1", path)
	if path == "/" {
		if t.root.isRealNode {
			return &resultFind{
				node: t.root,
			}
		}
		return nil
	}

	// return nil
	res := strings.Split(path, "/")
	depth := len(res)
	// depth := strings.Count(path, "/")

	// fmt.Println("depth", depth)

	// var nodesReal []*node
	// var queue []*node
	// queue = append(queue, t.root)

	// for len(queue) > 0 {
	// 	if queue[0].depth < depth {
	// 		// fmt.Println("current depth", queue[0].depth)
	// 		var queueTemp []*node
	// 		for _, n := range queue {
	// 			for _, childNode := range n.children {
	// 				queueTemp = append(queueTemp, childNode)
	// 			}
	// 		}
	// 		queue = queueTemp
	// 	} else {
	// 		// fmt.Println("2current depth", queue[0].depth)
	// 		// nodesReal = queue
	// 		break
	// 	}
	// }

	// for _, n := range nodesReal {
	// 	fmt.Println("waiting", n.fullpath, n.isRealNode)
	// }
	nodesReal, ok := t.allNodes[depth]
	if !ok {
		return nil
	}

	// fmt.Println("find", path, len(nodesReal))
	for _, n := range nodesReal {
		if n.isRealNode {
			// // fmt.Println(n.fullpath, n.fullPattern.exp.String())
			// if subMatch := n.fullPattern.exp.FindSubmatch([]byte(path)); subMatch != nil {
			// 	// matchParams := make(map[string]string)
			// 	// if strings.Index(n.fullpath, ":") > -1 {
			// 	// 	subMatch = subMatch[1:]
			// 	// 	// fmt.Println(n.fullPattern.exp.String(), n.fullPattern.keys, subMatch)
			// 	// 	for k, v := range subMatch {
			// 	// 		matchParams[n.fullPattern.keys[k]] = v
			// 	// 	}
			// 	// }

			// 	// return &resultFind{
			// 	// 	node:   n,
			// 	// 	params: matchParams,
			// 	// }
			// 	return nil
			// }
		}
	}
	// CHECK_NODE:
	// 	for _, n := range nodesReal {
	// 		if n.isRealNode {
	// 			// fmt.Println("")
	// 			// fmt.Println("check", n.fullpath)

	// 			nCurrent := n
	// 			var params = make(map[string]string)
	// 			for i := len(res) - 1; nCurrent != nil; i-- {
	// 				item := res[i]
	// 				// fmt.Println("check2", item, nCurrent.key, i, nCurrent.pattern)
	// 				if nCurrent.pattern != nil {
	// 					if !nCurrent.pattern.exp.MatchString(item) {
	// 						continue CHECK_NODE
	// 					}
	// 					params[nCurrent.pattern.name] = item
	// 				} else if i > 0 && nCurrent.key != item {
	// 					continue CHECK_NODE
	// 				}
	// 				nCurrent = nCurrent.parent
	// 			}
	// 			return &resultFind{
	// 				node:   n,
	// 				params: params,
	// 			}
	// 		}
	// 	}
	return nil
}

// Find returns nodes that the request match the route pattern
func (t *tree) Find(pattern string, isRegex bool) (nodes []*node) {
	var (
		node  = t.root
		queue []*Node
	)

	if pattern == node.path {
		nodes = append(nodes, node)
		return
	}

	if !isRegex {
		pattern = trimPathPrefix(pattern)
	}

	res := splitPattern(pattern)

	for _, key := range res {
		child, ok := node.children[key]

		if !ok && isRegex {
			break
		}

		if !ok && !isRegex {
			return
		}

		if pattern == child.path && !isRegex {
			nodes = append(nodes, child)
			return
		}
		node = child
	}

	queue = append(queue, node)

	for len(queue) > 0 {
		var queueTemp []*node
		for _, n := range queue {
			if n.isPattern {
				nodes = append(nodes, n)
			}

			for _, childNode := range n.children {
				queueTemp = append(queueTemp, childNode)
			}
		}

		queue = queueTemp
	}

	return
}

// trimPathPrefix is short for strings.TrimPrefix with param prefix `/`
func trimPathPrefix(pattern string) string {
	return strings.TrimPrefix(pattern, "/")
}

// splitPattern is short for strings.Split with param seq `/`
func splitPattern(pattern string) []string {
	return strings.Split(pattern, "/")
}
