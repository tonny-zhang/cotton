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
func TestTree1(t *testing.T) {
	tree := newTree()
	arrRouter := []string{
		"/",
		"/a",
		"/a/",
		"/a/:method",
		"/a/:method/:name",
		"/a/:method/:name/:id",
		"/a/:method/:name/:id/test",
		"/c/*file",
		"/b/",
		"/b/*file",
	}

	var tmpValue string
	for _, path := range arrRouter {
		tree.Add(path, func(p string) HandlerFunc {
			return func(c *Context) {
				tmpValue = "for " + p
				fmt.Println(tmpValue)
			}
		}(path))
	}

	// printTree(tree)

	result := tree.root.find("/a/vmethod")
	if result != nil {
		fmt.Println(result, result.node.key, result.params)
		result.node.handler(nil)
	} else {
		fmt.Println("no")
	}

	// assert.True(t, false)
}
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
		tree.Add(path, func(p string) HandlerFunc {
			return func(c *Context) {
				tmpValue = "for " + p
				// fmt.Println(tmpValue)
			}
		}(path))
	}
	tmpValue += ""

	// printTree(tree)

	assert.PanicsWithError(t, "[:test] in path [/a/:test] conflicts with [/a/:method]", func() {
		tree.Add("/a/:test", nil)
	})

	// 这里可能会随机得到冲突子元素
	assert.Panics(t, func() {
		tree.Add("/*file", nil)
	})
	assert.PanicsWithError(t, "action [*file] must end of path [/b/*file/test]", func() {
		tree.Add("/b/*file/test", nil)
	})
	assert.PanicsWithError(t, "action [*file] must end of path [/c/*file/test]", func() {
		tree.Add("/c/*file/test", nil)
	})
	assert.PanicsWithError(t, "path [test] must start with /", func() {
		tree.Add("test", nil)
	})

	// result := tree.Find("/a/vmethod")
	// if nil != result {
	// 	fmt.Println(result)
	// } else {
	// 	fmt.Println("no result")
	// }

	// assert.False(t, true)

	for _, path := range arrRouter {
		tmpValue = ""
		pathFind := strings.ReplaceAll(path, ":", "v")
		result := tree.Find(pathFind)
		// fmt.Println(path, pathFind, result)
		assert.NotNil(t, result)
		// fmt.Println("check", path, result)
		assert.NotNil(t, result.node.handler)
		result.node.handler(nil)
		assert.Equal(t, "for "+path, tmpValue)

		c := strings.Count(path, ":")
		if c > 0 {
			// fmt.Println("param  ->", path, result.params)
			assert.Equal(t, c, len(result.params), path)
		}
	}

	// result := tree.Find("/b/test/abc/123")
	// if nil != result {
	// 	result.node.handler(nil)
	// 	assert.Equal(t, "/test/abc/123", result.params["file"])
	// } else {
	// 	fmt.Println("no result")
	// }

	// t.Error("abc")
}
