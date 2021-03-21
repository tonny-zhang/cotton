[![GoDoc](https://pkg.go.dev/badge/github.com/tonny-zhang/cotton.svg)](https://pkg.go.dev/github.com/tonny-zhang/cotton) [![Release](https://img.shields.io/badge/release-v0.2.0-blue.svg?style=flat-square)](https://github.com/tonny-zhang/cotton/releases/tag/v0.2.0)

Cotton is a web framework written by Go (Golang).

# Contents
- [Contents](#contents)
	- [Installation](#installation)
	- [Quick start](#quick-start)
	- [Feature](#feature)
	- [API Example](#api-example)
		- [Using GET, POST, PUT, OPTIONS, DELETE, PATCH, HEAD](#using-get-post-put-options-delete-patch-head)
		- [Parameters in path](#parameters-in-path)
		- [Querystring parameters](#querystring-parameters)
		- [Using middleware](#using-middleware)
		- [Using group](#using-group)
	- [Benchmarks](#benchmarks)
	- [Author](#author)
	- [Acknowledgements](#acknowledgements)
## Installation
To install Cotton package, you need to install Go and set your Go workspace first.
1. The first need [Go](https://golang.org) installed
2. install Cotton
```sh
go get -u github.com/tonny-zhang/cotton
```
3. Import it in your code:
```go
import "github.com/tonny-zhang/cotton
```

## Quick start
You can find more in example/*

```go
package main

import "github.com/tonny-zhang/cotton"

func main() {
	r := cotton.NewRouter()
	r.Get("/hello", func(ctx *cotton.Context) {
		ctx.String("hello world from cotton")
	})

	r.Run(":8080")
}
```
## Feature
* Fast - see [Benchmarks](#benchmarks)
* router group
* parameters path
* middleware

## API Example
You can find a number of ready-to-run examples at [examples folder](./example)

### Using GET, POST, PUT, OPTIONS, DELETE, PATCH, HEAD

```go
func main() {
	r := cotton.NewRouter()
	r.Get("/hello", handler)
	r.Post("/hello", handler)

	r.Run(":8080")
}
```

### Parameters in path
```go
func main() {
	r := cotton.NewRouter()
	// /user/tonny		=> 	match
	// /user/123 		=> 	match
	// /user			=> 	no
	// /user/			=> 	no
	r.Get("/user/:name", func(c *cotton.Context) {
		c.String("hello "+c.Param("name"))
	})

	// /file/test		=> 	match
	// /file/a/b/c		=> 	match
	// /room/			=> 	no
	r.Get("/file/*file", func(c *cotton.Context) {
		c.String("file = "+c.Param("file"))
	})

	r.Run(":8080")
}
```

### Querystring parameters
```go
func main() {
	r := cotton.NewRouter()
	r.Get("/hello", func(c *cotton.Context) {
		name := c.GetQuery("name")
		first := c.GetDefaultQuery("first", "first default value")

		c.String("hello "+name+" "+first)
	})
	r.Run(":8080")
}
```

### Using middleware
```go
func main() {
	r := cotton.NewRouter()

	r.Use(cotton.Recover())
	r.Use(cotton.Logger())

	r.Get("/hello", func(c *cotton.Context) {
		c.String("hello")
	})
	r.Run(":8080")
}
```

### Using group
```go
func main() {
	r := cotton.NewRouter()
	g1 := r.Group("/v1", func(ctx *cotton.Context) {
		// use as a middleware in group
	})
	g1.Use(func(ctx *cotton.Context) {
		fmt.Println("g1 middleware 2")
	})
	{
		g1.Get("/a", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g1 a")
		})
	}

	r.Get("/v2/a", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "hello v2/a")
	})

	r.Run(":8080")
}
```

## Benchmarks
the benchmarks code for cotton be found in the [cotton-bench](https://github.com/tonny-zhang/cotton-bench) repository

## Author
* [tonny zhang](github.com/tonny-zhang)

## Acknowledgements
This package is inspired by the following
* [httprouter](github.com/xujiajun/gorouter)
* [gin](https://github.com/gin-gonic/gin)