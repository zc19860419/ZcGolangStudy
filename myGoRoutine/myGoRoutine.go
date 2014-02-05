package myGoRoutine

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var counter int = 0

func Count(lock *sync.Mutex) {
	lock.Lock()
	counter++
	fmt.Println(counter)
	lock.Unlock()
}
func testMutex() {
	lock := &sync.Mutex{}
	for i := 0; i < 10; i++ {
		go Count(lock)
	}
	for {
		lock.Lock()
		c := counter
		lock.Unlock()
		// when you switch on several chans, or (this is your case) when you explicitly
		// tell the scheduler to switch the contexts - this is what runtime.Gosched is for
		runtime.Gosched()
		if c >= 10 {
			break
		}
	}
}

func Count_chan(ch chan int) {
	ch <- 1
	fmt.Println(ch, "Counting")
}

func testChannel() {

	chs := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count_chan(chs[i])
	}
	for i, ch := range chs {
		<-ch               // This will block if the calculation is not done yet
		fmt.Println(ch, i) //0x20252300 0
	}
}

func testSelect() {
	ch := make(chan int, 1)
	for {
		select { //随机向ch中写入一个 0 或者 1
		case ch <- 1: //向 ch 写入数据 1
		case ch <- 0: //向 ch 写入数据 0
		}
		i := <-ch
		fmt.Println("Value received:", i)
	}
}

func testTimeOutAndDeadlock() {
	timeout := make(chan bool)
	ch := make(chan int)
	go func() {
		time.Sleep(1e9)
		timeout <- true
	}()

	go func() {
		time.Sleep(1e5)
		ch <- 1
	}()

	for {
		fmt.Println("for...")
		select {
		//等到两个匿名函数都执行了,由于2个 channel 都不带缓冲 ,so dead lock(blocking...)
		case i := <-ch:
			fmt.Println("Value received:", i)
		case <-timeout:
			fmt.Println("1sec time out")
		}
	}
}

type PipeData struct {
	value   int
	handler func(int) int
	next    chan int
}

//就可以达到流式处理数据的目的
func handle(queue chan *PipeData) {
	for data := range queue {
		data.next <- data.handler(data.value)
	}
}

func func1(value int) int {
	fmt.Println("func1")
	return value
}

func func2(value int) int {
	fmt.Println("func2")
	return value
}
func func3(value int) int {
	fmt.Println("func3")
	return value
}
func testPipe() {

	var (
		value1   = 1
		handler1 = func1
		c1       = make(chan int)
	)
	var data1 PipeData = PipeData{value1, handler1, c1}
	queue := make(chan *PipeData)
	go func() {
		queue <- &data1
		handle(queue)
	}()
	select {
	case data := <-queue:
		fmt.Println((*data).value)
	}
}

// func testPipe1() {

// 	var (
// 		value1   = 1
// 		handler1 = func1
// 		n1       = make(chan int)
// 	)
// 	var (
// 		value2   = 2
// 		handler2 = func2
// 		n2       = make(chan int)
// 	)
// 	var (
// 		value3   = 3
// 		handler3 = func3
// 		n3       = make(chan int)
// 	)
// 	var data1 PipeData = PipeData{value1, handler1, n1}
// 	var data2 PipeData = PipeData{value2, handler2, n2}
// 	var data3 PipeData = PipeData{value3, handler3, n3}
// 	datas := [3]PipeData{data1, data2, data3}
// 	queue := make(chan *PipeData)

// 	go func() {
// 		queue <- &datas
// 		handle(queue)
// 	}()

// 	for data := range datas {
// 		select {
// 		case <-(data.next):
// 			fmt.Println(data.value, "done")
// 		}
// 	}

// }

func testForwardChannel() {
	// var ch3 <-chan int // ch3是单向channel，只用于读取int数据

	//还支持类型转换
	ch4 := make(chan int)
	// ch5 := <-chan int(ch4) // ch5就是一个单向的读取channel
	// go func() {
	fmt.Println(0)
	ch4 <- 1
	// }()
	fmt.Println(2)
	i := <-ch4
	fmt.Println(i)
	// if _, ok := <-ch5; ok {
	// 	fmt.Println(ok)
	// } else {
	// 	fmt.Println(!ok)
	// }

	// defer close(ch3) // (cannot close receive-only channel)
	defer close(ch4)
	// defer close(ch5) // (cannot close receive-only channel)
}

func testBackwardChannel() {
	// var ch2 chan<- float64 // ch2是单向channel，只用于写float64数据

	// 还支持类型转换
	ch4 := make(chan int)
	// ch6 := chan<- int(ch4) // ch6 是一个单向的写入channel

	// defer close(ch2)
	defer close(ch4)
	// defer close(ch6)
}
func Test_goroutine() {
	// testMutex()

	// testChannel()
	// testSelect()

	// testTimeOutAndDeadlock()

	// testPipe()
	// testForwardChannel()
	// testBackwardChannel()
	testSyncLock()
}

// var l sync.Mutex
var test_any int

func testSyncLock() {
	i := 0
	ch := make(chan int /*, 10*/)
	for i < 30 {
		go func(int) {
			// l.Lock()
			// defer l.Unlock()
			fmt.Println(test_any, i)
			ch <- i
			test_any++
		}(i)
		k := <-ch
		fmt.Println("Received k:", k)
		i++
	}

	//下面这种方式不会报死锁错误,但是k始终读取出来时10
	// for i < 10 {
	// 	k := <-ch
	// 	fmt.Println("Received k:", k)
	// 	i++
	// }

	//下面这种方式会报死锁错误
	// for c := range ch {
	// 	fmt.Println("Received i:", c)
	// }
}
