# **gooh**
*gooh* is a web application framework written in [go](https://golang.org) with similar syntax to [express](http://expressjs.com) and [restify](http://restify.com) for building REST APIs on top of the go [http](https://golang.org/pkg/net/http) package.

## Getting it
You can get *gooh* using the following go command in the shell:
```golang
> go get github.com/goohland/gooh
```
once installed you can import it as follows:
```golang
import "github.com/goohland/gooh"
```

## Usage
A simple 'Hello World!' application in *gooh* looks like this:
```golang
// file main.go
package main

import (
	"github.com/goohland/gooh"
	"io"
	"net/http"
)

func main() {
	app := new(gooh.App)

	router := new(gooh.Router)
	router.GET("/hello", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
		io.WriteString(res, "Hello World!")
		return nil
	})
	app.AddMiddlewareHandler(router.GetMiddlewareHandler())

	http.HandleFunc("/", app.ServeHTTP)
	http.ListenAndServe(":8080", nil)
}
```
You can run it by executing in the shell:
```golang
> go run main.go
```
and test it in the browser by going to: [http://localhost:8080/hello](http://localhost:8080/hello)

## Middlewares
*gooh* defines a middleware as a function with the following type declaration:
```golang
type MiddlewareHandler func(*App, *Request, *Response) error
```
middleware handler functions will be called on every request in the order they were added to the *gooh* application:
```golang
// file main.go
package main

import (
	"github.com/goohland/gooh"
	"io"
	"net/http"
)

func HelloHandler(app *gooh.App, req *gooh.Request, res *gooh.Response) error {
	io.WriteString(res, "Hello")
	return nil
}

func WorldHandler(app *gooh.App, req *gooh.Request, res *gooh.Response) error {
	io.WriteString(res, " World!")
	return nil
}

func main() {
	app := new(gooh.App)
	app.AddMiddlewareHandler(HelloHandler)
	app.AddMiddlewareHandler(WorldHandler)

	http.HandleFunc("/", app.ServeHTTP)
	http.ListenAndServe(":8080", nil)
}
```
the code above will output 'Hello World!' too when you hit [http://localhost:8080](http://localhost:8080)


## Router
*gooh* comes with a built-in router, the router defines a `RouteHandler` as a function with the following type declaration:
```golang
type RouteHandler func(*App, *Request, *Response, map[string]string) error
```
you can add route handlers as follows:
```golang
router := new(gooh.Router)
router.AddRouteHandler("GET", "/hello", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
	io.WriteString(res, "Hello World!")
	return nil
})
```
the *gooh* router it's a middleware itself and you can add it to the *gooh* application the same way as any other middleware:
```golang
app.AddMiddlewareHandler(router.GetMiddlewareHandler())
```

### Convenience
*gooh* router offers five function wrappers (`GET`, `POST`, `PUT`, `DELETE` and `HEAD`) of the main `AddRouteHandler` function. The following two lines of code are equivalent:
```golang
router.GET("\users", gooh.Version{}, RouteHandler)
router.AddRouteHandler("GET", "\users", gooh.Version{}, RouteHandler)
```

### Versioning
*gooh* comes with built-in support for API [semantic versioning](http://semver.org), however *gooh* does not favor any particular versioning mechanism, instead it gives you the tools for you to implement versioning as you wish, the way it works is when you add a route handler to the router you need to specify the version
```golang
router.GET("/hello", gooh.Version{1,0,0}, RouteHandler)
```
and it exposes an `ApiVersion` property on the `gooh.Request` for you to set it to the correct version **before** getting to the router middleware
```golang
app.AddMiddlewareHandler(func(app *gooh.App, req *gooh.Request, res *gooh.Response) error {
	req.ApiVersion = &gooh.Version{1,0,0}
	return nil
})
...
app.AddMiddlewareHandler(router.GetMiddlewareHandler())
```
the router will then use this property to match against the registered route handlers, if a match is not found a `RouteNotFoundError` will be returned by the router middleware.

If you rather use url versioning simply specify your route path with the version and pass an empty `gooh.Version` to the router
```golang
router.GET("/v1/hello", gooh.Version{}, RouteHandler)
```

or if you do not want to version your API, you can remove the version from the url and pass an empty `gooh.Version` to the router
```golang
router.GET("/hello", gooh.Version{}, RouteHandler)
```
### Route Parameters
*gooh* offers built-in support for route parameters
```golang
router.GET("/users/:id", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
	io.WriteString(res, pms["id"])
	return nil
})
```
The code above will output the id sent as a parameter and match the given route handler against routes like [http://localhost:8080/users/123](http://localhost:8080/users/123) or  [http://localhost:8080/users/ABC](http://localhost:8080/users/ABC), but if you want to be more specific about the format of `:id` you can use a regular expression as follows:
```golang
router.GET("/users/:id{[0-9]+}", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
	io.WriteString(res, pms["id"])
	return nil
})
```
The code above will match the given route handler against routes like [http://localhost:8080/users/123](http://localhost:8080/users/123) but not against [http://localhost:8080/users/ABC](http://localhost:8080/users/ABC)

### Rules
*gooh* router enforces three rules and the router will panic if you try to break them

Only one route handler is allowed per route
```golang
router.GET("/hello", gooh.Version{}, HelloHandler)
router.GET("/hello", gooh.Version{}, WorldHandler)
```
```
panic: handler already exists for route: '/hello'
```

All parameters under the same route must have a unique name
```golang
router.GET("/users/:id/groups/:id", gooh.Version{}, RouteHandler)
```
```
panic: overwriting parameter: 'id' for route: '/users/:id/groups/:id'
```

If you provide a regular expression, it must be a valid one
```golang
router.GET("/users/:id{a)b}", gooh.Version{}, RouteHandler)
```
```
panic: error parsing regexp: unexpected ): `a)b`
```

## Context
*gooh* defines a `Context` interface as follows:
```golang
type Context interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
	Exists(string) (bool, error)
}
```
you can create your own implementations of the `Context` interface and use them to share data at two levels

### Application Level Context
Any data you want to initialize once and share during the lifespan of the application can be set in the application context as follows:
```golang
app := new(gooh.App)
app.Context = new(MyContextImplementation)
app.Context.Set("value", "hello")
```
and can be accessed from anywhere you have an instance of the `gooh.App` as follows:
```golang
router.GET("/hello", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
	value, _ := app.Context.Get("value")
	io.WriteString(res, value.(string))
	return nil
})
```

### Request Level Context
Any data you want to initialize on every request and share during the lifespan of the request can be set in the request context as follows:
```golang
app.AddMiddlewareHandler(func(app *gooh.App, req *gooh.Request, res *gooh.Response) error {
	req.Context = new(MyContextImplementation)
	req.Context.Set("value", "world")
	return nil
})
```
and can be accessed from anywhere you have an instance of the `gooh.Request` as follows:
```golang
router.GET("/hello", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
	value, _ := req.Context.Get("value")
	io.WriteString(res, value.(string))
	return nil
})
```

### Built-in Context
*gooh* offers a built-in `gooh.MemoryContext` which is nothing more than a `map[string]interface{}` implementing the `Context` interface, you can use it as follows:
```golang
app := new(gooh.App)
app.Context = new(gooh.MemoryContext)
app.Context.Set("value", 768)
val, _ := app.Context.Get("value")
fmt.Println(val)
```

## Error Handling
*gooh* defines an error handler as a function with the following type declaration:
```golang
type ErrorHandler func(*App, *Request, *Response, error)
```
if a middleware or a route handler returns an error, the chain of middlewares will be interrupted and the error handlers functions will be called in the order they were added to the *gooh* application:  
```golang
app.AddErrorHanlder(func(app *gooh.App, req *gooh.Request, res *gooh.Response, err error) {
	...
})
```
If an error is returned and no error handler was added, the application will panic with the error returned

If you route handler, middleware or an external service they call makes a panic call the *gooh* application will try to recover from it by calling the error handlers with a `gooh.PanicError` which contains the parameter sent to panic in the property `Err`

You can do type assertions in the error handler to deal with the different error types as follows:
```golang
app.AddErrorHanlder(func(app *gooh.App, req *gooh.Request, res *gooh.Response, err error) {
	switch err.(type) {
	case *gooh.RouteNotFoundError:
		http.NotFound(res, req.Request)
	case *gooh.PanicError:
		http.Error(res, err.Error(), 500)
		// or
		panic((err.(*gooh.PanicError)).Err)
	default:
		http.Error(res, err.Error(), 500)
	}
})
```

## Compatibility
The `gooh.Response` implements the `http.ResponseWriter` interface and the `gooh.Request` embeds the `http.Request` struct, which is why you can use any of the `http` package functions that take an `http.ResponseWriter` and/or an `http.Request` with the `gooh.Response` and `gooh.Request` as follows:
```golang
http.NotFound(res, req.Request)
```
in addition to the `http.ResponseWriter` methods, the `gooh.Response` exposes a `WriteJson` method to make it easy for json APIs, you can use it as follows:
```golang
router.GET("/hello", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
	res.WriteJson([]string{"Hello", "World!"})
	return nil
})
```

## Performance
*gooh* is intended to be a very thin layer on top of the http package, as such his performance it's almost identical to the http package, to compare it two hello world apps were created one in *gooh* and another one in the standard http package
```golang
// file with-gooh.go
package main

import (
	"github.com/goohland/gooh"
	"io"
	"net/http"
)

func main() {
	app := new(gooh.App)

	router := new(gooh.Router)
	router.GET("/hello/hello/hello/hello/hello/hello/hello/hello/hello/hello", gooh.Version{}, func(app *gooh.App, req *gooh.Request, res *gooh.Response, pms map[string]string) error {
		io.WriteString(res, "Hello World!")
		return nil
	})
	app.AddMiddlewareHandler(router.GetMiddlewareHandler())

	http.HandleFunc("/", app.ServeHTTP)
	http.ListenAndServe(":8080", nil)
}
```
```golang
// file without-gooh.go
package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/hello/hello/hello/hello/hello/hello/hello/hello/hello/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World")
	})
	http.ListenAndServe(":8888", nil)
}
```
> **Note:** the reason for the '/hello/hello/hello/hello/hello/hello/hello/hello/hello/hello' is so we can test the performance of the *gooh* router in deep routes

After running `ab -c 100 -n 1000 http://localhost:8080/hello/hello/hello/hello/hello/hello/hello/hello/hello/hello` multiple times you should get very similar results between the two applications like these:

```
// With gooh
Server Hostname:        localhost
Server Port:            8080

Document Path:          /hello/hello/hello/hello/hello/hello/hello/hello/hello/hello
Document Length:        12 bytes

Concurrency Level:      100
Time taken for tests:   0.178 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      129000 bytes
HTML transferred:       12000 bytes
Requests per second:    5612.30 [#/sec] (mean)
Time per request:       17.818 [ms] (mean)
Time per request:       0.178 [ms] (mean, across all concurrent requests)
Transfer rate:          707.02 [Kbytes/sec] received
```
```
// Without gooh
Server Hostname:        localhost
Server Port:            8080

Document Path:          /hello/hello/hello/hello/hello/hello/hello/hello/hello/hello
Document Length:        12 bytes

Concurrency Level:      100
Time taken for tests:   0.178 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      129000 bytes
HTML transferred:       12000 bytes
Requests per second:    5602.77 [#/sec] (mean)
Time per request:       17.848 [ms] (mean)
Time per request:       0.178 [ms] (mean, across all concurrent requests)
Transfer rate:          705.82 [Kbytes/sec] received
```

## License
The MIT License (MIT)
Copyright (c) 2012 Mark Cavage

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
