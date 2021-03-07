# Cotton Web Framework
Cotton is a web framework written by Go (Golang).

# Contents
- [Cotton Web Framework](#cotton-web-framework)
- [Contents](#contents)
	- [Installation](#installation)
	- [Quick start](#quick-start)
	- [Feature](#feature)
	- [API Example](#api-example)
		- [Using GET, POST, PUT, OPTIONS, DELETE, PATCH, HEAD](#using-get-post-put-options-delete-patch-head)
		- [Parameters in path](#parameters-in-path)
		- [Querystring parameters](#querystring-parameters)
		- [Using middleware](#using-middleware)
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
* router group
* regexp for router path
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

	// /room/123		=> 	match
	// /room			=> 	no
	// /room/			=> 	no
	// /room/tonny		=> 	no
	r.Get("/room/:id<num>", func(c *cotton.Context) {
		c.String("hello "+c.Param("id"))
	})

	// /action/123-ab		=> 	match
	// /action/1-aa			=> 	match
	// /action/11-bbb		=> 	no
	// /action/test			=> 	no
	r.Get("/action/:rule{\\d+-[ab]}", func(c *cotton.Context) {
		c.String("hello action "+c.Param("rule"))
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