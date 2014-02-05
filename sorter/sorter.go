package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sorter/algorithms/bubblesort"
	"sorter/algorithms/distributionCounting"
	"sorter/algorithms/insertSort"
	"sorter/algorithms/mergeSort"
	"sorter/algorithms/qsort"
	"strconv"
	"time"
)

var (
	infile    *string = flag.String("i", "infile", "File contains values for sorting")
	outfile   *string = flag.String("o", "outfile", "File to receive sorted values")
	algorithm *string = flag.String("a", "qsort", "Sort algorithm")
	verify    *bool   = flag.Bool("v", true, "Verify the result")
)

func main() {
	// 获取并解析命令行输入；
	fmt.Println("1.Parse commands")
	flag.Parse()
	if infile != nil {
		fmt.Println("infile =", *infile, "outfile =", *outfile, "algorithm =",
			*algorithm, "verify =", *verify)
	}
	// 从对应文件中逐行读取输入数据；
	fmt.Println("2.Read input")
	count, values, err := readValues(*infile)
	if err == nil {
		fmt.Println("Read ", count, " values:", values)
	} else {
		fmt.Println(err)
		return
	}

	// 调用对应的排序函数；
	fmt.Println("3.Sorting...")
	t1 := time.Now()

	switch *algorithm {
	case "qsort":
		qsort.QuickSort(values)
	case "bubblesort":
		bubblesort.BubbleSort(values)
	case "insert":
		insertSort.InsertionSort(values)
	case "merge":
		mergeSort.MergeSort(values)
	case "count":
		countingSort.CountingSort(values)
	default:
		fmt.Println("Sorting algorithm", *algorithm, "is either unknown or unsupported.")
	}

	t2 := time.Now()

	if *verify {
		fmt.Println("\t[optional]Verifying...")
		verifyValues(values)
	}

	// 将排序的结果输出到对应的文件中；
	fmt.Println("4.Output...")
	count, err = writeValues(values, *outfile)
	if err == nil {
		fmt.Println("Write", count, " values:", values)
	} else {
		fmt.Println(err)
	}
	// 打印排序所花费时间的信息。
	fmt.Println("4.Timing...")
	fmt.Println("The sorting process costs", t2.Sub(t1), "to complete.")
}

func verifyValues(values []int) {
	count := len(values)
	for i := 0; i < count-1; i++ {
		if values[i] > values[i+1] {
			fmt.Println("Panicking!")
			panic(fmt.Sprintf("values[%v] > values[%v]", i, i+1))
		}
	}
}

func readValues(infile string) (count int, values []int, err error) {
	fi, err := os.Open(infile)
	if err != nil {
		fmt.Println("Failed to open the input file ", infile)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	values = make([]int, 0)

	for {
		line, isPrefix, err1 := br.ReadLine()

		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		if isPrefix {
			fmt.Println("A too long line, seems unexpected.")
			return
		}

		str := string(line) // 转换字符数组为字符串

		value, err1 := strconv.Atoi(str) //转为数字

		if err1 != nil {
			err = err1
			return
		}

		values = append(values, value) //放入 int[] 数组中
		count++
	}
	return
}

func writeValues(values []int, outfile string) (count int, err error) {
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println("Failed to create the output file ", outfile)
		return 0, err
	}
	defer file.Close() //在返回 2 之前做的
	for _, value := range values {
		str := strconv.Itoa(value)
		file.WriteString(str + "\n")
		count++
	}
	return count, nil //2
}
