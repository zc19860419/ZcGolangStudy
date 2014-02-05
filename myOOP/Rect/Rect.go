package Rect

type Rect struct {
	x, y          float64
	width, height float64
}

func (r *Rect) Area() float64 {
	return r.height * r.width
}

func NewRect(x, y, width, height float64) *Rect {
	return &Rect{x, y, width, height}
}
