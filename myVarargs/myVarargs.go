package myVarargs

import (
	"fmt"
)

func myfunc1(args ...int) {

	for i, v := range args {
		fmt.Println("element[", i, "]=", v)
	}
}

//如果不指定类型 也是任意类型
func myPrintf(args ...interface{}) error {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Println(arg, "is an int value")
		case string:
			fmt.Println(arg, "is an string value")
		case int64:
			fmt.Println(arg, "is an int64 value")
		}
	}
	return nil
}

func Test_varargs() {
	var v1 int = 1
	var v2 int = 2
	var v3 int = 3
	var v4 int = 4

	myfunc1(v1, v2, v3, v4)

	var v5 int = 0
	var v6 string = "str"
	var v7 int64 = 7
	myPrintf(v5, v6, v7)
}
