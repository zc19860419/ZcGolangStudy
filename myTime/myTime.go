package myTime

import (
	"fmt"
	"time"
)

func Test_Mytime() {
	get_current_timestamp()
	// what_day_is_it_now()

	// show_before_after()

	var strings []string = []string{"abc", "efg", "doruimi"}
	for index, string := range strings {
		fmt.Println(index, ":", string)
	}
}

func show_before_after() {
	t := time.Date(2014, 1, 2, 0, 0, 0, 0, time.Local)
	b := time.Date(1999, 1, 1, 8, 0, 0, 0, time.Local)

	var str string
	if t.Before(b) {
		str = "before"
	} else {
		str = "after"
	}

	fmt.Println(t, str, b) //t的时间在b之前？  明显答案是false

	if t.After(b) {
		str = "after"
	} else {
		str = "before"
	}
	fmt.Println(t, str, b) //t的时间在b之后？   明显答案是true
}

func what_day_is_it_now() {
	//时间戳
	t := time.Now()
	fmt.Println(t.Weekday().String())
}

func get_current_timestamp() {
	//时间戳
	t := time.Now().Unix()
	fmt.Println(t)

	var test uint64 = 60
	var test1 uint64 /*= 18446744073709551615*/
	test1 = test
	fmt.Println(test1)

	timestamp := fmt.Sprintf("%d", t)
	fmt.Println("timestamp", timestamp)
	// 4294967296
	// 1390806728

	//时间戳到具体显示的转化
	fmt.Println(time.Unix(t, 0).String())

	//带纳秒的时间戳
	t = time.Now().UnixNano()
	fmt.Println(t)
	fmt.Println("------------------")

	//基本格式化的时间表示
	fmt.Println(time.Now().String())

	fmt.Println(time.Now().Format("2006year 01month 02day"))
}
