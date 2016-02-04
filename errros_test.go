package gooh

import (
	"errors"
	"testing"
)

func Test_RouteNotFoundError_Error(t *testing.T) {
	exp := "v"
	err := RouteNotFoundError{"v"}
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_PanicError_Error_String(t *testing.T) {
	exp := "v"
	err := PanicError{"v"}
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_PanicError_Error_Error(t *testing.T) {
	exp := "v"
	e := errors.New("v")
	err := PanicError{e}
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_PanicError_Error_Number(t *testing.T) {
	exp := "unknown panic error"
	err := PanicError{300}
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_PanicError_Error_Nil(t *testing.T) {
	exp := "unknown panic error"
	err := PanicError{nil}
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_PanicError_Error_Empty(t *testing.T) {
	exp := "unknown panic error"
	err := PanicError{}
	val := err.Error()

	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}
