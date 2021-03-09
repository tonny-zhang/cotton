package cotton

type (
	tree struct {
		root *node
	}
	node struct {
		key        string
		children   map[string]*node
		isRealNode bool
	}
)

func (t *tree) Add(path string, handler HandlerFunc) {
	var nodeCurrent = t.root
	var start, lenstr = 0, len(path)
Search:
	for {
		for i := start; i < lenstr; i++ {
			if path[i] == '/' {
				key := path[start:i]

				if _, ok := nodeCurrent.children[key]; !ok {

				}
				n := &node{}

				start = i + 1
				continue Search
			}
		}
		// items = append(items, Item{
		// 	key: str[start:],
		// })
		break
	}
}
