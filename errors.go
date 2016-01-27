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
	s := "Panic Error"
	if err, ok := e.Err.(error); ok {
		return s + ": " + err.Error()
	}
	return s
}
