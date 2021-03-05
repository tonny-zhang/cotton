package httpserver

type middleware struct {
	handlers    []HandlerFunc
	countRouter int
}
