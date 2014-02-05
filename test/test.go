package main

import (
	"fmt" //package implementing formatted I/O
	"math"
	"myAnonymous"
	"myClosure"
	"myErrorhandling"
	"myGoRoutine"
	"myMap"
	"myOOP"
	"myString"
	"myTime"
	"myVarargs"
	"slice"
	"strings"
)

func test_multi_assign() {
	var v1 string = "chenxi"
	var v2 string = "zhangchong"
	v1, v2 = v2, v1
	fmt.Printf(" %s\n", v2)
}

func test_multi_return() {
	_, lastname, _ := getName()
	fmt.Printf("lastname is %s\n", lastname)
}

func test_operator() {
	x := 2
	fmt.Printf("^2=%v\n", ^x)
}

func test_complex() {
	var value1, value2, value3 complex64 // 由2个float32构成的复数类型
	// value1 = 3.2 + 12i
	// value2 := 3.2 + 12i        // value2是complex128类型
	// value3 := complex(3.2, 12) // value3结果同value2

	value1, value2, value3 = 3.2+12i, 3.2+12i, complex(3.2, 12)

	fmt.Printf("value1=%v\n", value1)
	fmt.Printf("value2=%v\n", value2)
	fmt.Printf("value3=%v\n", value3)
}

func test_for(array [10]int) error {
	for i := 0; i < len(array); i++ {
		fmt.Println("Element", i, "of array is", array[i])
	}
	// Go语言还提供了一个关键字range，用于便捷地遍历容器中的元素。当然，数组也是range
	// 的支持范围。上面的遍历过程可以简化为如下的写法：
	for i, v := range array {
		fmt.Println("Array element[", i, "]=", v)
	}

	for pos, value := range "Go在中国" { //一个汉字占三个字节
		fmt.Printf("character '%c' type is %T value is %v, and start at byte position %d \n", value, value, value, pos)
		str := string(value) //convert rune to string
		fmt.Printf("string(%v)=>%s \n", value, str)
	}
	return nil

}

// p为用户自定义的比较精度，比如 0.00001
func IsEqual(f1, f2, p float64) bool {
	return true /*math.Fdim(f1, f2) < p*/
}

func getName() (firstName, lastName, nickName string) {
	if true {
		fmt.Printf("OK\n")
	}
	return "zhang", "chong", "handsome"

}

func test_rune() {
	var ch rune
	ch = 'e'
	fmt.Printf("ch v = %v\n", ch)
	fmt.Printf("ch c = %c\n", ch)
}

func test_string() {
	str := "Zhongguo"
	str1 := "Hello world" // 字符串也支持声明时进行初始化的做法
	// str[0] = 'X'        // 编译错误  字符串的内容不能在初始化后被修改
	total := str + " " + str1
	fmt.Printf("total = %s,total size is %v\n", total, len(total))

	str = "Hello,世界"
	n := len(str)
	for i := 0; i < n; i++ {
		ch := str[i] // 依据下标取字符串中的字符，类型为byte
		fmt.Println(i, ch)
	}
}

func test_banner(method_name string) {
	fmt.Printf("\n=======  " + method_name + "  ======\n")
}

func main() {
	// 指向数组的指针的声明。如：
	var p2array *[3]int
	// 这就是指针数组的声明。
	var pointers [3]*int

	f := func(x, y int) int {
		return x + y
	} //closable
	var v1 int = 10
	v1 := 11

	var v3 [10]int // 数组
	var v4 []int   // 数组切片
	var v5 struct {
		f int
	}
	var v6 *int           // 指针
	var v7 map[string]int // map，key为string类型，value为int类型
	var v8 func(a int) int
	fmt.Printf("%v", v1)

	test_banner("test_multi_assign")
	test_multi_assign()

	test_banner("test_multi_return")
	test_multi_return()

	test_banner("test_operator")
	test_operator()

	test_banner("test_complex")
	test_complex()

	test_banner("test_rune")
	test_rune()

	test_banner("test_string")
	test_string()

	test_banner("test_for")
	var array [10]int // int 数组如果没初始化 默认初始化为0
	test_for(array)

	test_banner("test_slice")
	slice.Test_slice()

	test_banner("test_varargs")
	myVarargs.Test_varargs()

	test_banner("Test_anonymous")
	myAnonymous.Test_anonymous()

	test_banner("Test_closure")
	myClosure.Test_closure()

	test_banner("Test_Errorhandling")
	myErrorhandling.Test_Errorhandling()

	test_banner("Test_Map")
	myMap.Test_Map()

	test_banner("Test_OOP")
	myOOP.Test_OOP()

	test_banner("Test_goroutine")
	myGoRoutine.Test_goroutine()

	testCopy()
	testRandom()

	test_banner("Test_String")
	myString.Test_String()

	test_banner("Test_Time")
	myTime.Test_Mytime()

	var f uint16 = 12
	fmt.Println(f, string(f))

}

func testRandom() {
	// Ni+1=(A* Ni + B)%  M          其中i = 0,1,...,M-1
	last := 23
	fmt.Println(last)
	for i := 0; i < 10000; i++ {
		next := (991*last + 857) % 10000
		fmt.Println(next)
		last = next
	}
}

func testCopy() {
	//                       0   1   2   3   4  5  6  7
	var values []int = []int{11, 22, 13, 4, 51, 6, 7, 8}
	left := 0
	right := 7
	m := 3
	i := left                   //[left ... i ... m]
	w := m + 1                  //[m+1  ... w ... right]
	s := values[left : right+1] //切片[start:end] end 是指超出末尾 1 个位置的那个地方,也就是实际是 [start,end-1] 的元素
	tmp := make([]int, 0)
	for i <= m && w <= right {
		if values[i] >= values[w] {
			tmp = append(tmp, values[w])
			fmt.Printf("%v,%v,append %v\n", values[i], values[w], values[w])
			w++
		} else {
			tmp = append(tmp, values[i])
			fmt.Printf("%v,%v,append %v\n", values[i], values[w], values[i])
			i++
		}
	}
	if i <= m {
		tmp = append(tmp, s[i:m+1]...)
		fmt.Printf("i <= m append s[%v:%v]\n", i, m+1)
	}

	if w <= right {
		tmp = append(tmp, s[w:right+1]...)
		fmt.Printf("w <= right append s[%v:%v]\n", w, right+1)
	}

	toString(tmp, "tmp")
	// 把临时数组的数据copy到原数组中
	copy(s[left:right+1], tmp[left:right+1])

	toString(s, "s")

	t := values[7:8]
	toString(t, "t")
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
