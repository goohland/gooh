package gooh

var (
	ErrRouteNotFound = &RouteNotFoundError{"route not found"}
)

type RouteNotFoundError struct {
	Msg string
}

func (e RouteNotFoundError) Error() string {
	return e.Msg
}

type PanicError struct {
	Err interface{}
}

func (e PanicError) Error() string {
	s := "unknown panic error"

	switch v := e.Err.(type) {
	case error:
		s = v.Error()
	case string:
		s = v
	}

	return s
}
