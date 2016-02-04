package gooh

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_Route_String_Value(t *testing.T) {
	exp := "v1 GET /users"
	route := Route{&Version{1, 0, 0}, "GET", "/users"}
	val := route.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Route_String_Empty(t *testing.T) {
	exp := ""
	route := Route{}
	val := route.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Route_String_EmptyVersion(t *testing.T) {
	exp := "GET /users"
	route := Route{&Version{}, "GET", "/users"}
	val := route.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Route_String_EmptyVersionAndMethod(t *testing.T) {
	exp := "/users"
	route := Route{&Version{}, "", "/users"}
	val := route.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Route_String_LowerMethod(t *testing.T) {
	exp := "GET /users"
	route := Route{&Version{}, "get", "/users"}
	val := route.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_Value(t *testing.T) {
	exp := []string{"users"}
	val := getPathFragments("/users")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_ValueWithoutSlahs(t *testing.T) {
	exp := []string{"users"}
	val := getPathFragments("users")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_ValuePrefixSlahses(t *testing.T) {
	exp := []string{"users"}
	val := getPathFragments("//users")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_ValueTrailingSlahs(t *testing.T) {
	exp := []string{"users"}
	val := getPathFragments("users/")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_ValueTrailingSlahses(t *testing.T) {
	exp := []string{"users"}
	val := getPathFragments("users//")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_Empty(t *testing.T) {
	exp := []string{""}
	val := getPathFragments("")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_getPathFragments_ValueMultipleSlahses(t *testing.T) {
	exp := []string{"users", "groups"}
	val := getPathFragments("/users/groups")

	if len(exp) != len(val) || strings.Join(val, "") != strings.Join(exp, "") {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_Nil(t *testing.T) {
	exp := "route not found"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{}, nil)
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_VersionNeverSet(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_Basic(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_Version(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{1, 7, 3}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 7, 3}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_WrongMethod(t *testing.T) {
	exp := "route not found"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "POST"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_LowerMethodRegistered(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("get", "/users", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_LowerMethodRequest(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "get"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_WrongPath(t *testing.T) {
	exp := "route not found"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/user"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_WrongVersion(t *testing.T) {
	exp := "route not found"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{1, 7, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{2, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_Regex(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/:id{[0-9]+}", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users/10"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_NotMatchingRegex(t *testing.T) {
	exp := "route not found"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/:id{[0-9]+}", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users/A"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_NestedPath(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/groups/friends", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users/groups/friends"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_AddRouteHandler_Parameters(t *testing.T) {
	exp := "uid:7;gid:10"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/:uid/groups/:gid", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("uid:" + pms["uid"] + ";gid:" + pms["gid"])
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{1, 0, 0}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users/7/groups/10"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_GET_Basic(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.GET("/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "GET"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_POST_Basic(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.POST("/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "POST"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_PUT_Basic(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.PUT("/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "PUT"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_DELETE_Basic(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.DELETE("/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "DELETE"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_HEAD_Basic(t *testing.T) {
	exp := "error"

	r := new(Router)
	r.HEAD("/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return errors.New("error")
	})
	mh := r.GetMiddlewareHandler()

	req := new(Request)
	req.Request = new(http.Request)
	req.Request.Method = "HEAD"
	req.ApiVersion = &Version{}
	req.Request.URL = new(url.URL)
	req.Request.URL.Path = "/users"

	err := mh(nil, req, nil)
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_String_Basic(t *testing.T) {
	exp := "GET /users"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return nil
	})
	val := r.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_String_Version(t *testing.T) {
	exp := "v1 GET /users"

	r := new(Router)
	r.AddRouteHandler("GET", "/users", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return nil
	})
	val := r.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_String_Parameter(t *testing.T) {
	exp := "v1 GET /users/:id"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/:id", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return nil
	})
	val := r.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_String_ParameterRegex(t *testing.T) {
	exp := "v1 GET /users/:id{[0-9]+}"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/:id{[0-9]+}", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return nil
	})
	val := r.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Router_String_NestedPath(t *testing.T) {
	exp := "v1 GET /users/:id{[0-9]+}/groups"

	r := new(Router)
	r.AddRouteHandler("GET", "/users/:id{[0-9]+}/groups", Version{1, 0, 0}, func(app *App, req *Request, res *Response, pms map[string]string) error {
		return nil
	})
	val := r.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}
