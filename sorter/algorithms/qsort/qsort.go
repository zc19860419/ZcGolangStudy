package qsort

func partition(values []int, low int, up int) int {
	pivot := values[up]
	i := low - 1
	for j := low; j < up; j++ {
		if values[j] <= pivot {
			i++
			values[i], values[j] = values[j], values[i]
		}
	}
	values[i+1], values[up] = values[up], values[i+1]
	return i + 1
}

func quickSort(values []int, low int, up int) {
	if low < up {
		mid := partition(values, low, up)
		//Watch out! The mid position is on the place, so we don't need to consider it again.
		//That's why below is mid-1, not mid! Otherwise it will occur overflow error!!!
		quickSort(values, low, mid-1)
		quickSort(values, mid+1, up)
	}
}

func QuickSort(values []int) {
	quickSort(values, 0, len(values)-1)
}
