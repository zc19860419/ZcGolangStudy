package myErrorhandling

import (
	"fmt"
	"io"
	"myErrorhandling/myError"
	"os"
)

func test_error() {
	// FIXME
	var a error = myError.New("stat", "MessageType is invalid!") //new(MyError)
	fmt.Println(a.Error())
	// if err.Err == nil {
	// 	fmt.Printf("err.Op=%s,err.Path=%s,err.Err = nil,ret = %v\n", err.Op, err.Path, ret)
	// } else {
	// 	fmt.Printf("err.Op=%s,err.Path=%s,err.Err = %s,ret = %v\n", err.Op, err.Path, err.Err.Error(), ret)
	// }

	// 	如果在处理错误时获取详细信息，而不仅仅满足于打印一句错误信息，那就需要用到类型转
	// 换知识了：err.(*os.PathError)
	_, err := os.Stat("a.txt")
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err != nil {
			// 获取PathError类型变量e中的其他信息并处理
			fmt.Println("err=>" + e.Error())
		}
	}
}

//只不过，当你需要为defer语句到底哪个先执行这种细节而烦恼的时候,
//说明你的代码架构可能需要调整一下了
func copyFile(dst, src string) (w int64, err error) {

	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println("err1=>" + err.Error())
		return
	}

	// defer func() { //<=the ith defer will (n-i+1)th call
	// 	fmt.Println("close src file")
	// 	close_count = close(srcFile, 0)() //undefined: close_count 变量的声明必须要在使用之前
	// 	close_count = close(srcFile, close_count)()
	// }()

	var close_count int
	close := func(file *os.File, count int) func() int {
		return func() int {
			if count == 0 {
				file.Close()
			}
			count++
			fmt.Println("close file (", count, ")times")
			return count
		}
	}

	defer func() { //<=the ith defer will (n-i+1)th call
		fmt.Println("close src file")
		close_count = close(srcFile, 0)()
		close_count = close(srcFile, close_count)()
	}()

	fmt.Println("os.Create(dst)...")
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("err2=>" + err.Error())
		return
	}

	defer func() { //<=the nth defer will 1th call
		fmt.Println("close dst File")
		close(dstFile, 0)()
	}()
	fmt.Println("io.Copy(dstFile, srcFile")
	return io.Copy(dstFile, srcFile)
}

func test_defer() {
	ok, err := copyFile("dst.txt", "src.txt")
	if err != nil {
		fmt.Printf("%v,err=>%s\n", ok, err.Error())
	}
}

// Go中可以抛出一个panic的异常，然后在 defer 中通过recover捕获这个异常，然后正常处理
func test_panic_recover() {
	// 必须要先声明defer，否则不能捕获到panic异常
	defer func() { //抛出一个panic的异常后,进入 defer 处理完以后程序直接退出,与 java 不同,不会继续往下走了!
		fmt.Println("c")
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		fmt.Println("d")
	}()
	fmt.Println("here1")
	fmt.Println("here2")
	// panic(404)
	panic("network broken")
	// panic(Error("file not exists"))
	fmt.Println("here3")
	// defer func() { // defer不能捕获到panic异常 因为在异常之前尚未声明
	// 	fmt.Println("e")
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println("f")
	// }()
	test_panic_recover()
}

func try(try_func func(string), exception_handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			exception_handler(err)
		}
	}()
	try_func("zhangchong")
}

type throw_panic func(string)

func real_panic(name string) {
	panic(name)
}

func test_try() {
	throws := real_panic
	try(throws, func(e interface{}) {
		print(e)
	})
}

type SsdbClient struct {
	host string
	port int
}

func (client *SsdbClient) CacheMessage(offline_client string) {
	fmt.Println("host:", client.host, ",port:", client.port, ".call", offline_client)
	panic(offline_client)
}
func (client *SsdbClient) try(try_func func(string), exception_handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			exception_handler(err)
		}
	}()
	try_func("zhangchong client")
}
func Test_Errorhandling() {

	// test_error()
	// test_defer()
	// test_panic_recover()
	// test_try()

	test_oop_try()
}

func test_oop_try() {
	client := new(SsdbClient)
	client.host = "www.baidu.com"
	client.port = 80

	throws := client.CacheMessage
	client.try(throws, func(e interface{}) {
		print(e)
	})
}
