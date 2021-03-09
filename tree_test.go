package cotton

import (
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
		"/a/:name",
		"/a/:method/:id/test",
		"/b/",
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
			// fmt.Println(path, result.params)
			assert.Equal(t, c, len(result.params))
		}
	}

	// result := tree.Find("/")
	// if nil != result {
	// 	result.node.handler(nil)
	// } else {
	// 	fmt.Println("no result")
	// }

	// t.Error("abc")
}
