package countingSort

import (
// "fmt"
)

const (
	MAX_RANGE = 100
)

func CountingSort(values []int) {
	N := len(values)
	//使用切片
	S := make([]int, N+1)
	COUNT := make([]int, MAX_RANGE)

	//D1[Clear COUNTs]
	for i := 0; i < MAX_RANGE; i++ {
		COUNT[i] = 0
	}

	//D2[Loop j]
	for j := 0; j < N; j++ {
		//D3[Increase COUNT[Kj]]
		COUNT[values[j]]++
	}

	//D4[Accumulate]
	for i := 1; i < MAX_RANGE; i++ {
		COUNT[i] = COUNT[i] + COUNT[i-1]
	}

	// for i := 1; i < MAX_RANGE; i++ {
	// 	fmt.Printf("COUNT[%d]=%d \n", i, COUNT[i])
	// }

	//D5[Loop on j](这时候 COUNT[i] 是小于等于 i 的键码的个数  特别地 COUNT[99]=N)
	for j := N - 1; j >= 0; j-- {
		//D6[Output Rj]
		i := COUNT[values[j]] //COUNT[values[j]] 是小于等于 values[j] 的键码的个数(i) 所以i即是最终排序后的所在位置
		S[i] = values[j]
		COUNT[values[j]] = i - 1
	}
	// 把临时数组的数据 copy到原数组切片中 这里切片是指针 是引用传递
	copy(values[:], S[1:])
}
