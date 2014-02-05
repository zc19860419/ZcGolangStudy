package bubblesort

func BubbleSort(values []int) /*[]int*/ { //根本无需返回一个数组 直接修改数组即可 这里传入的是引用
	flag := true
	for i := 0; i < len(values)-1; i++ {
		flag = true
		for j := 0; j < len(values)-i-1; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
				flag = false
			}
		} //end for j...
		if flag == true {
			break
		}
	} //end for i...
	/*return values*/
}
