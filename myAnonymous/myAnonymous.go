/**
Anonymous function
format:
func(parameters list) (return list){
	body:do something here
}(optional)=>if there is [Argument list],directly call the Anonymous function
*/

package myAnonymous

import (
	"fmt"
)

func test1() {
	callback := func(a, b int, z float64) bool {
		var (
			i int = 1
		)
		g_var++
		i++ //i无法持久化
		fmt.Println("global var =", g_var, ",local var i=", i)
		return a*b < int(z)
	}

	v1, v2 := 1, 2
	var v3 float64 = 1.0
	for i := 0; i < 3; i++ {
		fmt.Println("callback return ", callback(v1, v2, v3))
	}
}

var g_var int = 0

func test2() {
	v1, v2 := 1, 2
	var v3 float64 = 1.0
	var ret bool

	func(a, b int, z float64, ret *bool) {
		*ret = a*b < int(z)
	}(v1, v2, v3, &ret) //花括号后跟参数列表就是直接调用

	fmt.Println("Directly call anonymous return ", ret)

	ret = func(a, b int, z float64) bool {
		return a*b > int(z)
	}(v1, v2, v3) //花括号后跟参数列表就是直接调用
	fmt.Println("Directly call anonymous return ", ret)
}

func Test_anonymous() {
	test1()
	test2()
}
