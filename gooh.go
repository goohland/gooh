package gooh

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	numb := []int{v.Patch, v.Minor, v.Major}
	strs := []string{}
	valueFound := false

	for _, n := range numb {
		if n > 0 || valueFound {
			valueFound = true
			strs = append([]string{strconv.Itoa(n)}, strs...)
		}
	}
	var prefix string
	if len(strs) > 0 {
		prefix = "v"
	}
	return prefix + strings.Join(strs, ".")
}

func NewVersion(s string) *Version {
	v := &Version{}

	if matched, _ := regexp.MatchString("^v[0-9]+(\\.[0-9]+(\\.[0-9]+)?)?$", s); matched {
		numbers := strings.Split(strings.TrimPrefix(s, "v"), ".")
		for i, value := range numbers {
			number, _ := strconv.Atoi(value)
			switch i {
			case 0:
				v.Major = number
			case 1:
				v.Minor = number
			case 2:
				v.Patch = number
			}

		}
	}

	return v
}

type Context interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
	Exists(string) (bool, error)
}

type Request struct {
	*http.Request
	ApiVersion *Version
	Context    Context
}

type Response struct {
	http.ResponseWriter
}

func (r *Response) WriteJson(d interface{}) error {
	js, err := json.Marshal(d)
	if err != nil {
		return err
	}

	r.Header().Set("Content-Type", "application/json")
	r.Write(js)
	return nil
}

type MiddlewareHandler func(*App, *Request, *Response) error

type ErrorHandler func(*App, *Request, *Response, error)

type App struct {
	mdwHandlers []*MiddlewareHandler
	errHandlers []*ErrorHandler
	Name        string
	Version     *Version
	Context     Context
}

func (a *App) handleError(app *App, req *Request, res *Response, err error) {
	if len(a.errHandlers) == 0 {
		panic(err)
	}

	for _, handler := range a.errHandlers {
		(*handler)(app, req, res, err)
	}
}

func (a *App) AddMiddlewareHandler(h MiddlewareHandler) {
	if h != nil {
		a.mdwHandlers = append(a.mdwHandlers, &h)
	}
}

func (a *App) AddErrorHanlder(h ErrorHandler) {
	if h != nil {
		a.errHandlers = append(a.errHandlers, &h)
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &Request{r, &Version{}, nil}
	res := &Response{w}

	defer func() {
		if err := recover(); err != nil {
			a.handleError(a, req, res, &PanicError{err})
		}
	}()

	for _, handler := range a.mdwHandlers {
		if err := (*handler)(a, req, res); err != nil {
			a.handleError(a, req, res, err)
			return
		}
	}
}
