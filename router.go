package gooh

import (
	"regexp"
	"strings"
)

type RouteHandler func(*App, *Request, *Response, map[string]string) error

type Route struct {
	Version *Version
	Method  string
	Path    string
}

func (r Route) String() string {
	var v string
	if r.Version != nil {
		v = r.Version.String()
	}

	return strings.Trim(v+" "+strings.ToUpper(r.Method)+" "+r.Path, " ")
}

func getPathFragments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

type node struct {
	path     string
	param    string
	pattern  string
	handler  *RouteHandler
	children map[string]*node
}

func (n *node) addRouteHandler(r *string, f []string, p *map[string]bool, h *RouteHandler) {
	path := f[0]
	key := path
	var pattern string
	var param string

	if len(path) > 1 && strings.HasPrefix(path, ":") {
		key = "/"
		fIndex, lIndex := strings.Index(path, "{"), strings.LastIndex(path, "}")
		product := fIndex * lIndex
		switch {
		case product > 1:
			pattern = path[fIndex+1 : lIndex]
			if _, err := regexp.MatchString(pattern, ""); err != nil {
				panic(err)
			}
			param = path[1:fIndex]
			path = path[:fIndex]
		case product == 1:
			param = path[1:]
		case product < 0:
			panic("missing regex delimiter '{'' or '}' in: '" + path + "' for route: '" + (*r) + "'")
		}
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}

	if len(param) > 0 {
		if (*p)[param] == true {
			panic("overwriting parameter: '" + param + "' for route: '" + (*r) + "'")
		}
		(*p)[param] = true
	}

	child := n.children[key]
	if child == nil {
		child = new(node)
		child.path = path
		child.param = param
		child.pattern = pattern
		n.children[key] = child
	}

	if len(f) == 1 {
		if child.handler != nil {
			panic("handler already exists for route: '" + (*r) + "'")
		}
		child.handler = h
		return
	}

	child.addRouteHandler(r, f[1:], p, h)
}

func (n *node) buildRoutes(p string, r *[]string) {
	if len(n.children) == 0 {
		if n.handler != nil {
			*r = append(*r, p)
		}
		return
	}

	for _, child := range n.children {
		path := p + "/" + child.String()
		if child.handler != nil && len(child.children) > 0 {
			*r = append(*r, path)
		}
		child.buildRoutes(path, r)
	}
}

func (n node) String() string {
	pattern := n.pattern
	if len(pattern) > 0 {
		pattern = "{" + pattern + "}"
	}
	return n.path + pattern
}

type Router struct {
	trees map[string]*node
}

func (r *Router) addRouteHandler(method string, path string, v *Version, h *RouteHandler) {
	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = strings.TrimSuffix(path, "/")

	root := r.trees[v.String()]
	if root == nil {
		root = new(node)
		root.path = v.String()
		root.children = make(map[string]*node)
		r.trees[v.String()] = root
	}

	p := make(map[string]bool)
	root.addRouteHandler(&path, getPathFragments(strings.Join([]string{"/", strings.ToUpper(method), path}, "")), &p, h)
}

func (r *Router) getRouteHandler(method string, path string, v *Version) (*RouteHandler, map[string]string, error) {
	var rootKey string
	if v != nil {
		rootKey = v.String()
	}
	root := r.trees[rootKey]
	if root == nil {
		return nil, nil, ErrRouteNotFound
	}

	params := make(map[string]string)
	fragments := getPathFragments("/" + method + strings.TrimSuffix(path, "/"))
	var n *node
	for _, f := range fragments {
		n = root.children[f]

		if n == nil {
			n = root.children["/"]
			if n != nil {
				matched := len(n.pattern) == 0
				var err error

				if !matched {
					matched, err = regexp.MatchString(n.pattern, f)
					if err != nil {
						return nil, nil, err
					}
				}

				if matched {
					params[n.param] = f
				} else {
					n = nil
				}
			}
		}

		if n != nil {
			root = n
		} else {
			break
		}
	}

	if n == nil || n.handler == nil {
		return nil, nil, ErrRouteNotFound
	}

	return n.handler, params, nil
}

func (r *Router) AddRouteHandler(method string, path string, v Version, h RouteHandler) {
	if h != nil {
		r.addRouteHandler(method, path, &v, &h)
	}
}

func (r *Router) GET(path string, v Version, h RouteHandler) {
	r.addRouteHandler("GET", path, &v, &h)
}

func (r *Router) POST(path string, v Version, h RouteHandler) {
	r.addRouteHandler("POST", path, &v, &h)
}

func (r *Router) PUT(path string, v Version, h RouteHandler) {
	r.addRouteHandler("PUT", path, &v, &h)
}

func (r *Router) DELETE(path string, v Version, h RouteHandler) {
	r.addRouteHandler("DELETE", path, &v, &h)
}

func (r *Router) HEAD(path string, v Version, h RouteHandler) {
	r.addRouteHandler("HEAD", path, &v, &h)
}

func (r *Router) GetMiddlewareHandler() MiddlewareHandler {
	return func(app *App, req *Request, res *Response) error {
		if h, p, err := r.getRouteHandler(strings.ToUpper(req.Method), req.URL.Path, req.ApiVersion); err != nil {
			return err
		} else {
			return (*h)(app, req, res, p)
		}
	}
}

func (r *Router) GetRoutes() []*Route {
	routes := []*Route{}
	for _, node := range r.trees {
		version := NewVersion(node.String())

		for _, child := range node.children {
			method := child.String()

			paths := []string{}
			child.buildRoutes("", &paths)
			for _, path := range paths {
				routes = append(routes, &Route{version, method, path})
			}

		}
	}

	return routes
}

func (r Router) String() string {
	routes := r.GetRoutes()
	literals := []string{}

	for _, route := range routes {
		literals = append(literals, route.String())
	}

	return strings.Join(literals, "\n")
}
