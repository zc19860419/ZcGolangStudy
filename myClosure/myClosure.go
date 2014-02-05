/**
 * 闭包可以用来在一个函数与一组“私有”变量之间创建关联关系。
 * 在给定函数被多次调用的过程中，这些私有变量能够保持其持久性。
 * 变量的作用域仅限于包含它们的函数，因此无法从其它程序代码部分进行访问。
 * 不过，变量的生存期是可以很长，在一次函数调用期间所创建所生成的值在下次函数调用时仍然存在。
 * 正因为这一特点，闭包可以用来完成信息隐藏，并进而应用于需要状态表达的某些编程范型中。
 * 不过，用这种方式来使用闭包时，闭包不再具有引用透明性，因此也不再是纯函数。
 * 支持闭包的多数语言都将函数作为第一级对象，就是说这些函数可以存储到
 * 变量中作为参数传递给其他函数，最重要的是能够被函数动态创建和返回
 */
package myClosure

import (
	"fmt"
)

func test2() {
	var fn [10]func()

	for i := 0; i < len(fn); i++ {
		//这里的闭包其实就是一个形参列表和返回列表为另一个匿名函数的匿名函数的直接调用 在这个匿名函数中的私有变量可以保持其持久性
		fn[i] = func() func() {
			var k int = 0 //该私有变量可以持久化
			return func() {
				var j int = 0
				j++ //该私有变量无法持久化
				k++
				fmt.Println(i, k, j)
			}
		}()
	}
	for iter, f := range fn {
		fmt.Printf("======fn[%v]=======\n", iter)
		for i := 0; i < 3; i++ {
			fmt.Printf("\tcall f() %v time:", i)
			f()
		}
	}
}

func test1() {
	var y int = 5
	//==============================    [err start]    ==============================
	//cannot use func literal (type func()) as type func(...int) int in return argument
	// a := func() func(args ...int) int {
	// 	var i int = 10
	// 	return func() {
	// 		i++
	// 		j += 5
	// 		// 闭包包含着外部的环境变量值 j ,但这个环境变量值并不像匿名函数那样作为参数副本,
	// 		// 而是实实在在在的引用(或者指针,反正一个意思),当外部变量变化时,闭包能使用的值自然也就变化了
	// 		fmt.Printf("i, j: %d, %d\n", i, j)
	// 	}
	// }()
	// 闭包形式上就是一个返回值是另一个匿名函数的匿名函数的直接调用
	// Closure:[Anonymous function return_val.(type)=>sub Anonymous function
	//                    [sub Anonymous function{
	//                        return sub Anonymous function
	//                    }].directly_call
	//
	//==============================    [err end]    ==============================
	a := func() func(args ...int) int {
		var x int = 10
		return func(args ...int) int {
			x++
			fmt.Printf("x++ in anonymous func\n")
			y += 5
			// 闭包包含着外部的环境变量值 j ,但这个环境变量值并不像匿名函数那样作为参数副本,
			// 而是实实在在在的引用(或者指针,反正一个意思),当外部变量变化时,闭包能使用的值自然也就变化了
			fmt.Printf("x, y: %d, %d\n", x, y)
			err := 0
			return err
		}
	}()
	a()
	fmt.Println("next call a()")
	y *= 2
	a()
}

/**
 * 我的理解：1.闭包函数是把创建时，引用到的外部数据复制了一份，与函数一起组成了一个整体。
 * 闭包函数出现的条件：
 * 1.被嵌套的函数引用到非本函数的外部数据，而且这外部数据不是“全局变量”
 * 2.函数被独立了出来(被父函数返回或赋值给其它函数或变量了)
 *
 * 还找到句明言：对象是附有行为的数据，而闭包是附有数据的行为。
 */
func ExFunc1(n int) func() {
	sum := n
	a := func() { //把匿名函数作为值赋给变量a (Go 不允许函数嵌套。然而你可以利用匿名函数实现函数嵌套)
		fmt.Println(sum + 1) //调用本函数外的变量
	} //这里没有()匿名函数不会马上执行
	return a
}

//这个函数还有另一个写法：
func ExFunc2(n int) func() {
	sum := n
	return func() { //直接在返回处的匿名函数
		fmt.Println(sum + 1)
	}
}

func testExFunc() {
	myFunc := ExFunc1(10)
	myFunc()
	myAnotherFunc := ExFunc1(20)
	myAnotherFunc()

	myFunc()
	myAnotherFunc()
}

func ExFunc3(n int) func() {
	sum := n
	a := func() {
		sum++ //在这里对外部数据加1
		fmt.Println(sum)
	}
	return a
}

func testExFunc1() {
	myFunc1 := ExFunc3(10)
	myFunc1()
	myAnotherFunc1 := ExFunc3(20)
	myAnotherFunc1()
	myFunc1()
	myAnotherFunc1()
	//这里得出的结果是22，由此可以证明两点
	//1.闭包中对外部数据的修改，外部不可见
	//2.外部数据的值被保存到新建的静态变量中
}

/**
*The first way is a trivial example using variable f and function myInc(),
but this shows the function type inc being declared then a matching function (it has one int input parameter and returns an int)
and a variable f that is assigned the function myInc.
第一个方法是一个使用变量f和myInc()函数的实验性例子 , Println(2+1) 将给出相同的结果
只是这样展示了函数类型 inc 的声明然后一个匹配的函数 myInc 和一个被赋值的变量 f

Go uses the Duck typing method to see if a function matches.
Duck typing simply means if it looks like a duck and sounds like a duck then it is treated as if it were a duck.
In coding terms, if the function has the same input and return parameters as those specified for the type
then the compiler is happy and allows it to be used.
Go 使用了"鸭子" 类型方法来检查一个函数是否匹配.
鸭子类型简而言之:只要它像鸭子一样走路，像鸭子一样嘎嘎叫，那它就是只鸭子.
在吧编程术语中,如果函数与某种指定类型的函数有相同的输入和返回参数,那么编译器就将乐于允许他被使用.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
在面向对象设计思想中，有这样一个重要原则：对接口编码，不对实现编码①。如果利用鸭子类型，实现这一原则只需极少的额外工
作，轻轻松松就能完成。举个例子，对象若有push和pop方法，它就能当作栈来用；反之若没有，就不能当作栈。
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
If you've used function pointers in C/C++ then you'll have an idea about assigning a function to a variable.
Golang is treating functions as first class functions , so they can be passed as functions to other functions and even returned.
The second way though is considerably more complicated.
The function getIncbynFunction returns an anonymous function that when it is called
will add the value passed into the original call. Let's take it step by step.
getIncbynFunction将返回一个匿名函数,该匿名函数被调用时将加上一个传入父函数调用的参数.让我们一步一步实现它

g := getIncbynFunction
This assigns the function getIncbynFunction to g, so we can call g with an int parameter and it will
return an anonymous function. So
这将 getIncbynFunction 赋值给 g,所以我们传入一个 int 参数来调用 g,然后我们将返回一个匿名函数,所以

h :=g(4)
This assigns the anonymous function returned by g (aka getIncbynFunction)
with the anonymous function return statement effectively written as return value + 4.
So then calling h(5) returns 9. The second assignment
这就将 g(也就是 getIncbynFunction) 返回的匿名函数赋值给了 h,同时匿名函数的返回语句被有效地写成了 "return value + 4"
所以之后调用  h(5) 返回 9

i :=g(6)
Hooks up variable i to the anonymous function returned by g but with its return statement as return value + 6.
Calling i(1) returns 7. This is pretty whizzy and takes a bit of getting your head round.
将变量 i 与 由 g 返回的匿名函数挂上钩,就是将返回语句写成了 "value + 6".
调用 i(1) 返回 7. 这太棒了而且需要多动动你的脑子

It's a Closure!这就是闭包!
This is what makes function types different from a C function pointer, not to mention easier to understand syntax
and much better type safety. The h= g(4) not only binds the getIncbynFunction to h,
but it includes its environment so you get access to the non local variable- in this case the variable n.
That's what makes it possible to call the anonymous function via two different variables h and i and have different values of n.
There's a considerably more sophisticated example of anonymous functions and closures in the Pig Game .
It has functions that use other functions as arguments and return values.
它使得函数类型与 C 函数指针不一样了,更不要说理解其语法会更容易和更好的类型安全性.不仅将 h与 getIncbynFunction 绑定了,还包括了
他的环境上下文所以只有你就能够访问那些非局部变量了-比如这里的 n.
这就使得通过两个不同的变量h和i 以及任意不同的 n 的值 来调用匿名函数 成为了可能.
在Pig Game 中有一个匿名函数和闭包的复杂得多的例子.
它有使用其他函数作为参数和返回值的函数。

Closures come out of functional programming and are also found in C#.
*/
type inc func(digit int) int

func myInc(value int) int {
	return value + 1
}

func getIncbynFunction(n int) inc {
	var test int = 10
	return func(value int) int {
		test++
		fmt.Println("getIncbynFunction", test)
		return value + n
	}
}

func testGetIncbynFunction() {
	f := myInc
	g := getIncbynFunction
	h := g(4)
	i := g(6)
	fmt.Println(f(2))
	fmt.Println(h(5))
	fmt.Println(i(1))
}

// 原理
// –闭包只是带有父函数的上下文(Context)的函数
// • 函数带有上下文并不奇怪，函数都可以访问全局变量，那就是上下文。
// • 不同之处在于，父函数本身的状态是动态产生和消亡的，这个上下文需要有生命周期管理。但 Go 是
//   gc 语言，这一点上也不是问题。
// • 闭包对父函数的Context只是引用而不复制。
// 闭包 (closure)
// • 柯里化 (currying)
//     –对多元函数的某个参数进行绑定
//     	func app(in io.Reader, out io.Writer, args []string) { ... }
//      	args := []string{"arg1", "arg2", ...}
//         	app2 := func(in io.Reader, out io.Writer) {
//         	app(in, out, args)
//     	}
//     	–Go1.1 支持了对 receiver 的快速绑定
//     	func (recvr *App) main(in io.Reader, out io.Writer) { ... }
//     	app := &App{...}
// 	   	app2 := app.main
// • 等价于
// 		app2 := func(in io.Reader, out io.Writer) { app.main(in, out) }
// 闭包 (closure)
// • 闭包的组合
// 		func pipe(
// 			app1 func(io.Reader, io.Writer),
// 			app2 func(io.Reader, io.Writer)
// 		) func(io.Reader, io.Writer) {
// 			return func(in io.Reader, out io.Writer) {
// 				pr, pw := io.Pipe()
// 				defer pw.Close()
// 				go func() {
// 				defer pr.Close()
// 				app2(pr, out)
// 		}()
// 			app1(in, pw)
// 		}
// 		}
// 闭包 (closure)
// • 闭包的陷阱
// 		var closures [2]func()
// 		for i := 0; i < 2; i++ {
// 			closures[i] = func() {
// 				fmt.Println(i)
// 			}
// 		}
// 		closures[0]()
// 		closures[1]()
// – 不要认为会打印 0 和 1，实际打印是 2 和 2
// – 修正方法
// 		var closures [2]func()
// 			for i := 0; i < 2; i++ {
// 				val := i
// 			closures[i] = func() {
// 				fmt.Println(val)
// 			}
// 		}
// 		closures[0]()
// 		closures[1]()
func Test_closure() {
	test1()
	test2()
	testExFunc()
	testExFunc1()
	testGetIncbynFunction()
}
