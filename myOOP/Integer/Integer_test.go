package Integer

import (
	"testing"
)

func TestAdd(t *testing.T) {
	var a Integer = 1
	var b Integer = 2
	var c *Integer = &a
	c.Add(b)
	if *c != 3 {
		t.Error("Integer Less() failed.Got ", *c, "Expected 3")
	}
}

func TestLess(t *testing.T) {
	var a Integer = 1
	var b Integer = 2
	if a.Less(b) == false {
		t.Error("Integer Less() failed.Got false ,Expected true")
	}
}
