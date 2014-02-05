package main

import (
	"fmt" //package implementing formatted I/O
	// "error"
	// "log"
)

const (
	total_cond   = 1048576
	question_num = 10
)

func main() {
	solution()
}

//使用枚举法来完成
func solution() {
	// logfile, err := os.OpenFile("./solution_log", os.O_RDWR|os.O_CREATE, 0)
	// if err != nil {
	// 	fmt.Printf("%s\r\n", err.Error())
	// 	os.Exit(-1)
	// }
	// defer logfile.Close()
	// logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	var option int32 /*= 129176*/
	for option = 0; option < total_cond; option++ {
		options := parse(option)
		// toString(options, false)
		// fmt.Printf("\n")
		errno, result := verify(options)
		if false == result {
			// fmt.Println("Verify", errno, "failed")
			continue
		} else {
			fmt.Printf("**********************    The solution is:%v    **********************\n", errno)
			toString(options, true)
			fmt.Printf("**********************    The solution end    **********************\n")
		}
	}
	fmt.Println(option)
	//等待控制台输入
	// running := true
	// reader := bufio.NewReader(os.Stdin)
	// for running {
	// 	data, _, _ := reader.ReadLine()
	// 	command := string(data)
	// 	if command == "stop" {
	// 		running = false
	// 	}
	// 	log.Println("command", command)
	// }
}

func parse(option int32) (options [question_num + 1]uint) {
	var i uint
	tmp := uint(option)
	for i = 1; i <= question_num; i++ {
		options[i] = (tmp & 0x3) + 1
		tmp = (tmp >> 2)
	}
	return options
}

func verify(options [question_num + 1]uint) (uint, bool) {
	if false == verify_q1(options) {
		return 1, false
	}
	if false == verify_q2(options) {
		return 2, false
	}
	if false == verify_q3(options) {
		return 3, false
	}
	if false == verify_q4(options) {
		return 4, false
	}
	if false == verify_q5(options) {
		return 5, false
	}
	if false == verify_q6(options) {
		return 6, false
	}
	if false == verify_q7(options) {
		return 7, false
	}
	if false == verify_q8(options) {
		return 8, false
	}
	if false == verify_q9(options) {
		return 9, false
	}
	if false == verify_q10(options) {
		return 10, false
	}
	return 0, true
}

func toString(options [question_num + 1]uint, newline bool) {

	for i, v := range options {
		if i == 0 {
			continue
		}
		if newline {
			switch v {
			case 1:
				fmt.Printf("%v:A\n", i)
			case 2:
				fmt.Printf("%v:B\n", i)
			case 3:
				fmt.Printf("%v:C\n", i)
			case 4:
				fmt.Printf("%v:D\n", i)
			default:
				panic("invalid options")
			}
		} else {
			switch v {
			case 1:
				fmt.Printf("%v:A,", i)
			case 2:
				fmt.Printf("%v:B,", i)
			case 3:
				fmt.Printf("%v:C,", i)
			case 4:
				fmt.Printf("%v:D,", i)
			default:
				panic("invalid options")
			}
		}

	}
}

//the first A is A:1,B:2,C:3,D:4
func verify_q1(options [question_num + 1]uint) bool {

	switch {
	case options[1] == 1:
		return true
	case options[2] == 1 && options[1] != 1:
		return true
	case options[3] == 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[4] == 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[5] == 1 && options[4] != 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[6] == 1 && options[5] != 1 && options[4] != 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[7] == 1 && options[6] != 1 && options[5] != 1 && options[4] != 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[8] == 1 && options[7] != 1 && options[6] != 1 && options[5] != 1 && options[4] != 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[9] == 1 && options[8] != 1 && options[7] != 1 && options[6] != 1 && options[5] != 1 && options[4] != 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	case options[10] == 1 && options[9] != 1 && options[8] != 1 && options[7] != 1 && options[6] != 1 && options[5] != 1 && options[4] != 1 && options[3] != 1 && options[2] != 1 && options[1] != 1:
		return true
	}
	return false
}

func Unique(options [question_num + 1]uint, min int) (int, bool) {
	for i := 1; i <= question_num-1; i++ {
		j := i + 1
		if options[i] == options[j] && i != min {
			// fmt.Println(i, "==", j)
			return i, false
		}

	}
	return 0, true
}

//"唯一" "连续两个"具有相同答案的问题是 A:5_6,B:6_7,C:7_8,D:8_9
func verify_q2(options [question_num + 1]uint) bool {
	var unique bool

	switch options[2] {
	case 1:
		_, unique = Unique(options, 5)
		if options[5] == options[6] && unique {
			return true
		}
	case 2:
		_, unique = Unique(options, 6)
		if options[6] == options[7] && unique {
			return true
		}
	case 3:
		_, unique = Unique(options, 7)
		if options[7] == options[8] && unique {
			return true
		}
	case 4:

		_, unique = Unique(options, 8)
		if options[8] == options[9] && unique {
			return true
		}
	}
	// fmt.Println("2=>")
	return false
}

//第三个问题和哪一个问题答案相同 A:4,B:9,C:8,D:2
func verify_q3(options [question_num + 1]uint) bool {
	switch options[3] {
	case 1:
		if options[3] == options[4] {
			return true
		}
	case 2:
		if options[3] == options[9] {
			return true
		}
	case 3:
		if options[3] == options[8] {
			return true
		}
	case 4:
		if options[3] == options[2] {
			return true
		}
	}
	return false
}

func countABCD(options [question_num + 1]uint) (int, int, int, int) {
	var (
		countA int
		countB int
		countD int
		countC int
	)
	for i := 1; i <= question_num; i++ {
		if options[i] == 1 {
			countA++
		}
		if options[i] == 2 {
			countB++
		}
		if options[i] == 3 {
			countC++
		}
		if options[i] == 4 {
			countD++
		}
	}
	return countA, countB, countC, countD
}

// the #A is A:5,B:4,C:3,D:2
func verify_q4(options [question_num + 1]uint) bool {
	numA, _, _, _ := countABCD(options)
	switch options[4] {
	case 1:
		if numA == 5 {
			return true
		}
	case 2:
		if numA == 4 {
			return true
		}
	case 3:
		if numA == 3 {
			return true
		}
	case 4:
		if numA == 2 {
			return true
		}
	}
	return false
}

//第五个问题和哪一个问题答案相同 A:1,B:2,C:3,D:4
func verify_q5(options [question_num + 1]uint) bool {
	switch options[5] {
	case 1:
		if options[5] == options[1] {
			return true
		}
	case 2:
		if options[5] == options[2] {
			return true
		}
	case 3:
		if options[5] == options[3] {
			return true
		}
	case 4:
		if options[5] == options[4] {
			return true
		}
	}
	return false
}

//#A=?#X A:nil,B:C,C:C,D:D
func verify_q6(options [question_num + 1]uint) bool {
	numA, numB, numC, numD := countABCD(options)
	switch options[6] {
	case 1:
		if numA != numB && numA != numC && numA != numD {
			return true
		}
	case 2:
		fallthrough
	case 3:
		if numA == numC {
			return true
		}
	case 4:
		if numA == numD {
			return true
		}
	}
	return false
}

//本问题答案与下一题相差 A:3,B:2,C:1,D:0
func verify_q7(options [question_num + 1]uint) bool {
	switch {
	case options[7] == 1 && options[8] == 4:
		return true
	case options[7] == 2 && options[8] == 4:
		return true
	case options[7] == 3 && (options[8] == 2 || options[8] == 4):
		return true
	case options[7] == 4 && options[8] == 4:
		return true
	}

	return false

}

// the #A is A:0,B:1,C:2,D:3
func verify_q8(options [question_num + 1]uint) bool {
	numA, _, _, _ := countABCD(options)
	switch options[8] {
	case 1:
		if numA == 0 {
			return true
		}
	case 2:
		if numA == 1 {
			return true
		}
	case 3:
		if numA == 2 {
			return true
		}
	case 4:
		if numA == 3 {
			return true
		}
	}
	return false
}

func isCompositeNum(n int) bool {
	if n == 0 || n == 1 || n == 4 || n == 6 || n == 8 || n == 9 || n == 10 {
		return true
	}
	return false
}
func isPrime(n int) bool {
	if n == 2 || n == 3 || n == 5 || n == 7 {
		return true
	}
	return false

}
func isSquare(n int) bool {
	if n == 4 || n == 9 {
		return true
	}
	return false
}

// the #!A is A:合数,B:质数,C:<5,D:平方
func verify_q9(options [question_num + 1]uint) bool {
	_, numB, numC, numD := countABCD(options)
	numConsonant := numB + numC + numD

	switch options[9] {
	case 1:
		if isCompositeNum(numConsonant) {
			return true
		}
	case 2:
		if isPrime(numConsonant) {
			return true
		}
	case 3:
		if numConsonant < 5 {
			return true
		}
	case 4:
		if isSquare(numConsonant) {
			return true
		}
	}
	// fmt.Println("9=>")
	// fmt.Println(numConsonant)
	return false
}

// 本题答案 is A:A,B:B,C:C,D:D
func verify_q10(options [question_num + 1]uint) bool {
	return true
}
