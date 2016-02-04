package gooh

import (
	"testing"
)

func Test_MemoryContext_Exists_False(t *testing.T) {
	exp := false
	c := new(MemoryContext)

	val, _ := c.Exists("k")
	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_MemoryContext_Exists_True(t *testing.T) {
	exp := true
	c := new(MemoryContext)

	c.Set("k", "v")
	val, _ := c.Exists("k")
	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_MemoryContext_Get_Unset(t *testing.T) {
	var exp interface{}
	c := new(MemoryContext)

	val, _ := c.Get("k")
	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_MemoryContext_Get_Set_Value(t *testing.T) {
	exp := "v"
	c := new(MemoryContext)

	c.Set("k", "v")
	val, _ := c.Get("k")
	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}

func Test_MemoryContext_Get_Set_Nil(t *testing.T) {
	var exp interface{}
	c := new(MemoryContext)

	c.Set("k", nil)
	val, _ := c.Get("k")
	if val != exp {
		t.Errorf("Expected '%v', got '%v'", exp, val)
	}
}
