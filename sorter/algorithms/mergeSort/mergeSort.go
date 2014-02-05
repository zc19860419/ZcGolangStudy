package mergeSort

import (
	"fmt"
)

func MergeSort(values []int) {
	goroutine := false
	mergeSort(values, 0, len(values)-1, goroutine)
}

//切片[start:end] end 是指超出末尾 1 个位置的那个地方,也就是实际是 [start,end-1] 的元素
func merge(values []int, left, m, right int) {
	// fmt.Printf("merge (%v %v %v)...\n", left, m, right)
	// fmt.Println(values[left:m+1], values[m+1:right+1])
	i := left  //[left ... i ... m]         slice:[left:m+1]
	w := m + 1 //[m+1  ... w ... right-1]   slice:[m+1:right+1]
	tmp := make([]int, 0)
	for i <= m && w <= right {
		if values[i] >= values[w] {
			tmp = append(tmp, values[w])
			w++
		} else {
			tmp = append(tmp, values[i])
			i++
		}
	}
	if i <= m {
		tmp = append(tmp, values[i:m+1]...)
	}

	if w <= right {
		tmp = append(tmp, values[w:right+1]...)
	}

	// 把临时数组的数据copy到原数组切片中,这里切片是指针 是引用传递
	copy(values[left:right+1], tmp[0:])
	// toString(values, "values")
}

//在Go语言的官方FAQ中描述, maps / slices / channels 是引用类型, 数组是值类型
func mergeSort(values []int, left, right int, goroutine bool) {
	// fmt.Printf("mergeSort (%v %v)\n", left, right)
	if left >= right {
		// fmt.Printf("mergeSort return\n")
		return
	}

	m := (left + right) / 2

	if !goroutine {
		mergeSort(values, left, m, goroutine)
		mergeSort(values, m+1, right, goroutine)
		merge(values, left, m, right)
	} else {
		c1 := make(chan int)
		c2 := make(chan int)

		defer func() {
			close(c1)
			close(c2)
		}()
		go func() {
			mergeSort(values, left, m, goroutine)
			c1 <- 1
		}()
		go func() {
			mergeSort(values, m+1, right, goroutine)
			c2 <- 2
		}()

		<-c1
		<-c2
		merge(values, left, m, right)
	}
}

func toString(s []int, name string) {
	fmt.Printf("%s:", name)
	for i, v := range s {
		if i > 0 {
			fmt.Printf(",")
		}
		fmt.Printf("%v", v)
	}
	fmt.Printf("\n")
}
