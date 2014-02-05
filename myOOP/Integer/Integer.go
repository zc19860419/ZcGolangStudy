package Integer

type Integer int

func (a Integer) Less(b Integer) bool {
	return a < b
}

// func (a *Integer) Add(b Integer) Integer {
// 	ret := (*a + b)
// 	return ret
// }

func (a *Integer) Add(b Integer) {
	*a += b
}
