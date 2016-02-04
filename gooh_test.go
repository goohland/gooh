package gooh

import (
	"errors"
	"testing"
)

var gvar int

func Test_Version_String_Value(t *testing.T) {
	exp := "v99.99.99"
	v := Version{99, 99, 99}
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_String_NegativeValue(t *testing.T) {
	exp := ""
	v := Version{-99, -99, -99}
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_String_MissingPatchValue(t *testing.T) {
	exp := "v1.2"
	v := Version{}
	v.Major = 1
	v.Minor = 2
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_String_MissingMinorValue(t *testing.T) {
	exp := "v1.0.7"
	v := Version{}
	v.Major = 1
	v.Patch = 7
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_String_MissingPatchAndMinorValue(t *testing.T) {
	exp := "v1"
	v := Version{}
	v.Major = 1
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_String_MissingMajorrValue(t *testing.T) {
	exp := "v0.1.7"
	v := Version{}
	v.Minor = 1
	v.Patch = 7
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_String_Empty(t *testing.T) {
	exp := ""
	v := Version{}
	val := v.String()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_Empty(t *testing.T) {
	exp := new(Version)
	val := NewVersion("")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_Value(t *testing.T) {
	exp := &Version{99, 99, 99}
	val := NewVersion("v99.99.99")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_InvalidValue(t *testing.T) {
	exp := &Version{}
	val := NewVersion("v-99.-99.-99")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_MissingPatchValue(t *testing.T) {
	exp := &Version{}
	exp.Major = 1
	exp.Minor = 2
	val := NewVersion("v1.2")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_MissingPatchAndMinorValue(t *testing.T) {
	exp := &Version{}
	exp.Major = 1
	val := NewVersion("v1")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_ZeroPatchValue(t *testing.T) {
	exp := &Version{}
	exp.Major = 1
	exp.Minor = 2
	exp.Patch = 0
	val := NewVersion("v1.2.0")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_ZeroMinorValue(t *testing.T) {
	exp := &Version{}
	exp.Major = 1
	exp.Patch = 2
	val := NewVersion("v1.0.2")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_ZeroPatchAndMinorValue(t *testing.T) {
	exp := &Version{}
	exp.Major = 1
	exp.Minor = 0
	exp.Patch = 0
	val := NewVersion("v1.0.0")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_Version_New_ZeroMajorValue(t *testing.T) {
	exp := &Version{}
	exp.Major = 0
	exp.Minor = 50
	exp.Patch = 22
	val := NewVersion("v0.50.22")

	if *val != *exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_App_AddMiddlewareHandler_Invalid(t *testing.T) {
	exp := 0
	app := new(App)
	app.AddMiddlewareHandler(nil)
	val := len(app.mdwHandlers)

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_App_AddMiddlewareHandler_Nil(t *testing.T) {
	var exp interface{}
	app := new(App)
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		return nil
	})
	val := (*(app.mdwHandlers[0]))(nil, nil, nil)

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_App_AddMiddlewareHandler_Error(t *testing.T) {
	exp := errors.New("error")
	app := new(App)
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		return errors.New("error")
	})
	val := (*(app.mdwHandlers[0]))(nil, nil, nil)

	if val.Error() != exp.Error() {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_App_AddErrorHandler_Invalid(t *testing.T) {
	exp := 0
	app := new(App)
	app.AddErrorHanlder(nil)
	val := len(app.errHandlers)

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_App_AddErrorHandler_Fn(t *testing.T) {
	gvar = 1
	exp := 7
	app := new(App)
	app.AddErrorHanlder(func(app *App, req *Request, res *Response, err error) { gvar = 7 })
	(*(app.errHandlers[0]))(nil, nil, nil, nil)

	if gvar != exp {
		t.Errorf("Expected '%v', got '%v'", exp, gvar)
	}
}

func Test_App_handleError_FnAndOrder(t *testing.T) {
	gvar = 0
	exp := 3
	app := new(App)
	app.AddErrorHanlder(func(app *App, req *Request, res *Response, err error) { gvar += 7 })
	app.AddErrorHanlder(func(app *App, req *Request, res *Response, err error) { gvar -= 4 })
	app.handleError(nil, nil, nil, nil)

	if gvar != exp {
		t.Errorf("Expected '%v', got '%v'", exp, gvar)
	}
}

func Test_App_ServeHTTP_FnAndOrder(t *testing.T) {
	gvar = 0
	exp := 3
	app := new(App)
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		gvar += 7
		return nil
	})
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		gvar -= 4
		return nil
	})
	app.ServeHTTP(nil, nil)

	if gvar != exp {
		t.Errorf("Expected '%v', got '%v'", exp, gvar)
	}
}

func Test_App_ServeHTTP_Error(t *testing.T) {
	gvar = 0
	exp := 3
	app := new(App)
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		return errors.New("error")
	})
	app.AddErrorHanlder(func(app *App, req *Request, res *Response, err error) {
		if err.Error() == "error" {
			gvar = 3
		}
	})
	app.ServeHTTP(nil, nil)

	if gvar != exp {
		t.Errorf("Expected '%v', got '%v'", exp, gvar)
	}
}

func Test_App_ServeHTTP_PanicWithString(t *testing.T) {
	gvar = 0
	exp := 3
	app := new(App)
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		panic("error")
	})
	app.AddErrorHanlder(func(app *App, req *Request, res *Response, err error) {
		if e, ok := err.(*PanicError); ok && e.Error() == "error" {
			gvar = 3
		}
	})
	app.ServeHTTP(nil, nil)

	if gvar != exp {
		t.Errorf("Expected '%v', got '%v'", exp, gvar)
	}
}

func Test_App_ServeHTTP_PanicWithError(t *testing.T) {
	gvar = 0
	exp := 3
	app := new(App)
	app.AddMiddlewareHandler(func(app *App, req *Request, res *Response) error {
		panic(errors.New("error"))
	})
	app.AddErrorHanlder(func(app *App, req *Request, res *Response, err error) {
		if e, ok := err.(*PanicError); ok && e.Error() == "error" {
			gvar = 3
		}
	})
	app.ServeHTTP(nil, nil)

	if gvar != exp {
		t.Errorf("Expected '%v', got '%v'", exp, gvar)
	}
}
