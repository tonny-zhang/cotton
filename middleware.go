package cotton

type middleware struct {
	handlers    []HandlerFunc
	countRouter int
}
