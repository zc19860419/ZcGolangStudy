package myMap

import (
	"fmt"
)

const (
	REGISTER = iota
	SUBSCRIBE
)

type Type struct {

	// const (//结构体中不允许有枚举型
	// 	REGISTER = iota
	// 	SUBSCRIBE
	// )
}

// SubInfo 是一个包含订阅消息信息的结构
type SubInfo struct {
	AppKey    string
	AppValue  string
	TopicType int
	AppName   string
	// numberOfFields int //not exported
}

var subDB map[string]SubInfo

func test_insert() {
	//元素赋值
	subDB["1"] = SubInfo{"1234343", "Jack", REGISTER, "Room 101,..."}

	// 往这个map里插入几条数据
	subDB["12345"] = SubInfo{"1234wqw5", "Tom", SUBSCRIBE, "Room 203,..."}
	subDB["12"] = SubInfo{"sd1", "Lucy", SUBSCRIBE, "Room 511,..."}
}

func test_search() {
	// 从这个map查找键为"1234"的信息
	sub, ok := subDB["1234"]

	// ok是一个返回的bool型，返回true表示找到了对应的数据
	if ok {
		fmt.Println("Found subInfo", sub.AppName, "with ID 1234.")
	} else {
		fmt.Println("Did not find subInfo with ID 1234.")
	}

	value, ok := subDB["1234343"]
	if ok { // 找到了
		// 处理找到的value
		fmt.Println("Found subInfo", value, "with ID 1234.")
	}

}

func test_traverse() {
	for key, value := range subDB {
		fmt.Println("key:", key, ", value:{AppKey=", value.AppKey, ",AppName=",
			value.AppName, ",AppValue=", value.AppValue, ",TopicType=", value.TopicType, "}")
	}
}

func test_make() {
	subDB = make(map[string]SubInfo) //string是键的类型， SubInfo 则是其中所存放的值类型
	//我们可以使用Go语言内置的函数make()来创建一个新map
	//也可以选择是否在创建时指定该map的初始存储能力
	// myMap := make(map[string]SubInfo, 100) //下面的例子创建了一个初始存储能力为100的map:
}

func test_delete() {
	//元素删除
	delete(subDB, "1")
	fmt.Println("delete key:1")
}

func Test_Map() {

	test_make()
	test_insert()
	test_search()
	test_traverse()
	test_delete()
	test_traverse()
}
