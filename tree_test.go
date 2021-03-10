package cotton

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	tree := newTree()
	arrRouter := []string{
		"/",
		"/a",
		"/a/",
		"/a/:id",
		"/a/:action/:name",
		"/a/:method/:id/test",
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

	// tree.print()
	assert.PanicsWithError(t, "[/a/:test] conflicts with [/a/:id]", func() {
		tree.Add("/a/:test", nil)
	})
	assert.PanicsWithError(t, "path [/*file] conflicts with other rule", func() {
		tree.Add("/*file", nil)
	})
	assert.PanicsWithError(t, "[/b/*file/test] conflicts with [/b/*file]", func() {
		tree.Add("/b/*file/test", nil)
	})
	assert.PanicsWithError(t, "action [*file] must end of path [/c/*file/test]", func() {
		tree.Add("/c/*file/test", nil)
	})
	assert.PanicsWithError(t, "path [test] must start with /", func() {
		tree.Add("test", nil)
	})

	// assert.False(t, true)

	for _, path := range arrRouter {
		tmpValue = ""
		result := tree.Find(strings.ReplaceAll(path, ":", "v"))
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

	result := tree.Find("/b//test/abc/123")
	if nil != result {
		result.node.handler(nil)
		assert.Equal(t, "/test/abc/123", result.params["file"])
	} else {
		fmt.Println("no result")
	}

	// t.Error("abc")
}
