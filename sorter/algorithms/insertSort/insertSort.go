package insertSort

import (
	"container/list"
)

func InsertionSort(values []int) {
	straightInsert(values)
	// shellInsert(values)
}

var h_aux []int = []int{1, 2, 4, 8}

func shellInsert(values []int) {
	N := len(values)
	t := len(h_aux)
	//D1[Loop on s]
	s := t - 1
	for ; s >= 0; s-- {
		//D2 [Loop on j]
		h := h_aux[s]            //h=8
		for j := h; j < N; j++ { //j=8
			//D3
			i := j - h //i=0
			K := values[j]
			for i >= 0 {
				//D4 if K >= values[i] ,back to D6
				if K >= values[i] {
					break
				}
				//D5 if i>0 ,back to D4
				values[i+h] = values[i]
				i = i - h
			}
			//D6[R into R(i+h)]
			values[i+h] = K // (ignore in h-order...)...<=values[j-2h]<=values[j-h]<=values[j] (ignore in h-order...)
		}
	}
}

func straightInsert(values []int) {
	sortedData := list.New()
	sortedData.PushBack(values[0])
	size := len(values)
	for i := 1; i < size; i++ {
		v := values[i]
		e := sortedData.Front()
		for nil != e {
			if e.Value.(int) >= v {
				sortedData.InsertBefore(v, e)
				// sortedData.InsertAfter(v, e)
				break
			}
			e = e.Next()
		}
		//the biggest,put @v on the back of the list
		if nil == e {
			sortedData.PushBack(v)
		}
	}

	//变量还是有作用域 这里要重新声明
	i := 0
	e := sortedData.Front()
	for nil != e {
		values[i] = e.Value.(int)
		e = e.Next()
		i++
	}
}
