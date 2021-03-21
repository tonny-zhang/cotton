package cotton

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func printTree(tree *tree) {
	fmt.Println("=============")
	tree.root.print(0)
	fmt.Println("=============")
}
func TestAddConflicts(t *testing.T) {
	assert.PanicsWithError(t, "path [test] must start with /", func() {
		tree := newTree()
		tree.add("test", nil)
	})
	assert.PanicsWithError(t, "[:method] in path [/a/:method] conflicts with [/a/:test]", func() {
		tree := newTree()
		tree.add("/a/:test", nil)
		tree.add("/a/:method", nil)
	})

	assert.PanicsWithError(t, "[:method] in path [/a/:method] conflicts with [/a/*file]", func() {
		tree := newTree()
		tree.add("/a/*file", nil)
		tree.add("/a/:method", nil)
	})

	assert.PanicsWithError(t, "[*file] in path [/a/*file] conflicts with [/a/:method]", func() {
		tree := newTree()
		tree.add("/a/:method", nil)
		tree.add("/a/*file", nil)
	})

	assert.PanicsWithError(t, "action [*file] must end of path [/a/*file/test]", func() {
		tree := newTree()
		tree.add("/a/*file/test", nil)
	})

	assert.PanicsWithError(t, "path [/:test/12] conflicts with [/:test/:id]", func() {
		tree := newTree()
		tree.add("/:test/:id", nil)
		tree.add("/:test/12", nil)
	})

	assert.NotPanics(t, func() {
		tree := newTree()
		tree.add("/:test/:id", nil)
		tree.add("/:test/12/abc", nil)
	})

	assert.PanicsWithError(t, "path [/:test/:id] conflicts with [/:test/12]", func() {
		tree := newTree()
		tree.add("/:test/12", nil)
		tree.add("/:test/:id", nil)
	})
	assert.PanicsWithError(t, "path [/:test/12] conflicts with [/:test/*file]", func() {
		tree := newTree()
		tree.add("/:test/*file", nil)
		tree.add("/:test/12", nil)
	})
	assert.PanicsWithError(t, "path [/:test/*file] conflicts with [/:test/12]", func() {
		tree := newTree()
		tree.add("/:test/12", nil)
		tree.add("/:test/*file", nil)
	})
}

// func TestTree1(t *testing.T) {
// 	tree := newTree()
// 	arrRouter := []string{
// 		"/",
// 		"/a",
// 		"/a/",
// 		"/a/:method",
// 		"/a/:method/:name",
// 		"/a/:method/:name/:id",
// 		"/a/:method/:name/:id/test",
// 		"/c/*file",
// 		"/b/",
// 		"/b/*file",
// 	}

// 	var tmpValue string
// 	for _, path := range arrRouter {
// 		tree.add(path, func(p string) HandlerFunc {
// 			return func(c *Context) {
// 				tmpValue = "for " + p
// 				// fmt.Println(tmpValue)
// 			}
// 		}(path))
// 	}
// 	tmpValue += ""

// 	printTree(tree)

// 	var pathTest string

// 	pathTest = "/a/test/abc/123/test"
// 	pathTest = "/c/test/abc/123/test"
// 	pathTest = "/b/"
// 	result := tree.root.find(pathTest)
// 	if result.node != nil {
// 		result.node.handler(nil)
// 		fmt.Println(result, result.node.key, result.params, tmpValue)
// 	} else {
// 		fmt.Println("no")
// 	}

// 	assert.True(t, false)
// }
func TestTree(t *testing.T) {
	tree := newTree()
	arrRouter := []string{
		"/",
		"/a",
		"/a/",
		"/a/:method",
		"/a/:method/:name",
		"/a/:method/:name/test",
		"/b/",
		"/b/*file",
	}

	var tmpValue string
	for _, path := range arrRouter {
		tree.add(path, func(p string) HandlerFunc {
			return func(c *Context) {
				tmpValue = "for " + p
				// fmt.Println(tmpValue)
			}
		}(path))
	}
	tmpValue += ""

	// printTree(tree)

	for _, path := range arrRouter {
		tmpValue = ""
		pathFind := strings.ReplaceAll(path, ":", "v")
		result := tree.root.find(pathFind)
		assert.NotNil(t, result)
		assert.NotNil(t, result.node.handler)
		result.node.handler(nil)
		assert.Equal(t, "for "+path, tmpValue)

		c := strings.Count(path, ":")
		if c > 0 {
			assert.Equal(t, c, len(result.params), path)
		}
	}
}
