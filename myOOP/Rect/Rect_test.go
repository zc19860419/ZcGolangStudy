package Rect

import (
	"testing"
)

func TestArea(t *testing.T) {
	// r := &Rect{1.0, 2.0,
	// 	3.0, 4.0}
	var (
		x, y          float64 = 1.0, 2.0
		width, height float64 = 3.0, 4.0
	)
	r := NewRect(x, y, width, height)
	area := r.Area()
	if area != 12.0 {
		t.Error("Rect Area() failed.Got ", area, "Expected 12.0")
	}
}
