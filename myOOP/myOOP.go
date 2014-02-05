package myOOP

import (
	"fmt"
	"myOOP/Integer"
)

func testLess() {
	var a Integer.Integer = 1
	var b Integer.Integer = 2
	if a.Less(b) {
		fmt.Printf("(%v) less than (%v)\n", a, b)
	}
}

func testAdd() {
	var a Integer.Integer = 1
	var b Integer.Integer = 2
	var c *Integer.Integer = &a
	c.Add(b)
}

//Go语言中的数组和基本类型没有区别，是很纯粹的值类型
func testValue() {
	var a = [3]int{1, 2, 3}
	var b = a
	b[1]++
	fmt.Println(a, b)
}

// 其一，Go语言的标准库，再也不需要绘制类库的继承树图。你一定见过不少C++、Java、C#
// 类库的继承树图。这里给个Java继承树图：
// http://docs.oracle.com/javase/1.4.2/docs/api/overview-tree.html
// 在Go中，类的继承树并无意义，你只需要知道这个类实现了哪些方法，每个方法是啥含义
// 就足够了。
// 其二，实现类的时候，只需要关心自己应该提供哪些方法，不用再纠结接口需要拆得多细才
// 合理。接口由使用方按需定义，而不用事前规划。
// 其三，不用为了实现一个接口而导入一个包，因为多引用一个外部的包，就意味着更多的耦
// 合。接口由使用方按自身需求来定义，使用方无需关心是否有其他模块定义过类似的接口
type Engine interface {
	Start()
	Stop()
}

type Starter interface {
	Start()
}

type Stoper interface {
	Stop()
}

type MyEngine struct {
}

func (m *MyEngine) Start() {
	fmt.Printf("MyEngine Start...\n")
}
func (m *MyEngine) Stop() {
	fmt.Printf("MyEngine Stop...\n")
}

type Car struct {
	Engine
}

func (c *Car) Drive() {
	fmt.Printf("Getting in the car ...\n")
	c.Start()
	c.Stop()
	fmt.Printf("Getting out of the car ...\n")
}

func NewCar(engine Engine) *Car {
	return &Car{engine}
}

func testAnonymousCombination() {
	var myengine Engine = new(MyEngine)
	mycar := NewCar(myengine)
	mycar.Drive()

	//接口查询
	if starter, ok := myengine.(Starter); ok {
		fmt.Printf("Engine=>Starter,starter.Start...\n")
		starter.Start()
	}
}

const (
	STARTER     = 1 << iota
	STOPER      = 1 << iota
	STARTSTOPER = STARTER | STOPER
)

func gettype(v1 interface{}) (what int) {
	isStarter := false
	switch v1.(type) {
	case Starter:
		isStarter = true
		what = STARTER
	case Stoper:
		what = STOPER
		if isStarter {
			what = STARTSTOPER
		}
	default:
		panic("unrecognized Interface type")
	}
	return
	// switch v := v1.(type) {
	// case int: // 现在v的类型是int
	// case string: // 现在v的类型是string
	// default:
	// 	if v, ok := arg.(Stringer); ok { // 现在v的类型是Stringer
	// 		val := v.String()
	// 		// ...
	// 	} else {
	// 		// ...
	// 	}
	// }
}

func testQueryInterface() {
	var myengine Engine = new(MyEngine)
	var v1 interface{} = myengine

	what := gettype(v1)
	fmt.Println("what is ", what)
	fmt.Println("STARTER is ", STARTER)
	fmt.Println("STOPER is ", STOPER)
	fmt.Println("STARTSTOPER is ", STARTSTOPER)
	switch what { //swicth case 不能用于 type 的选择上 所以这里使用了 gettype 方法

	case STARTER:
		fmt.Printf("I am Starter...\n")
		var v = v1.(Starter)
		v.Start()
		fallthrough //fallthrough 是无论如何都会进入下一个 case! 才不管 case 的条件判断
	case STOPER:
		fmt.Printf("I am Stoper...\n")
		var v = v1.(Stoper)
		v.Stop()
	case STARTSTOPER:
		fmt.Printf("I am Stoper & Stoper...\n")
		var v = v1.(Engine)
		v.Start()
		v.Stop()
	default:
		panic("unrecognized Interface type")
	}

	// Unpack 4 bytes into uint32 to repack into base 85 5-byte.
	// var v uint32
	// switch len(src) {
	// default:
	// 	v |= uint32(src[3])
	// 	fallthrough
	// case 3:
	// 	v |= uint32(src[2]) << 8
	// 	fallthrough
	// case 2:
	// 	v |= uint32(src[1]) << 16
	// 	fallthrough
	// case 1:
	// 	v |= uint32(src[0]) << 24
	// }
}

func Test_OOP() {
	// testLess()
	// testAdd()
	// testValue()
	// testAnonymousCombination()
	testQueryInterface()
}
