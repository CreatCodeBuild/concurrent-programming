package main

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

func main() {
	input := randomList(10000 * 10000)
	timeIt(func() {
		result3 := mergeSortParallel3(input)
		fmt.Println(len(result3))
	})
	timeIt(func() {
		result3 := mergeSortParallel5(input)
		fmt.Println(len(result3))
	})
	// timeIt(func() {
	// 	result := mergeSort(input)
	// 	fmt.Println(len(result))
	// })	
	// // fmt.Println(result)
	// timeIt(func() {
	// 	result2 := mergeSortParallel(input)
	// 	fmt.Println(len(result2))
	// })
	// fmt.Println(result2)
	
}

// merge sort
func mergeSort(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	left, right := split(list)
	left = mergeSort(left)
	right = mergeSort(right)
	return merge(left, right)
}

func split(list []int) ([]int, []int){
	mid := len(list) / 2
	return list[0:mid], list[mid:]
}

func merge(left, right []int) []int {
	result := make([]int, len(left) + len(right))
	l := 0
	r := 0
	for ;l+r < len(result) ; {
		if l < len(left) && r < len(right) {
			if left[l] < right[r] {
				result[l+r] = left[l]
				l += 1
			} else {
				result[l+r] = right[r]
				r += 1
			}
		} else if l >= len(left) {
			result[l+r] =  right[r]
			r += 1
		} else {
			result[l+r] =  left[l]
			l += 1
		}
	}
	return result
}

func mergeSortParallel(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	left, right := split(list)
	c := make(chan []int)
	go func() {
		c <- mergeSort(left)
	}()
	go func() {
		c <- mergeSort(right)
	}()
	return merge(<- c, <- c)
}

func mergeSortParallel2(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	left, right := split(list)
	c := make(chan []int)
	go func() {
		c <- mergeSortParallel2(left)
	}()
	go func() {
		c <- mergeSortParallel2(right)
	}()
	return merge(<- c, <- c)
}

func mergeSortParallel3(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	left, right := split(list)
	c := make(chan []int)
	limit := 3125000 * 4
	go func() {
		if len(left) <= limit {
			c <- mergeSort(left)
		} else {
			c <- mergeSortParallel3(left)
		}
	}()
	go func() {
		if len(right) <= limit {
			c <- mergeSort(right)
		} else {
			c <- mergeSortParallel3(right)
		}
	}()
	return merge(<- c, <- c)
}

func mergeSortParallel4(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	left, right := split(list)
	c := make(chan []int)
	limit := 16
	go func() {
		if runtime.NumGoroutine() >= limit {
			c <- mergeSort(left)
		} else {
			c <- mergeSortParallel4(left)
		}
	}()
	go func() {
		if runtime.NumGoroutine() >= limit {
			c <- mergeSort(right)
		} else {
			c <- mergeSortParallel4(right)
		}
	}()
	return merge(<- c, <- c)
}

func mergeSortParallel5(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	left, right := split(list)
	limit := 3125000 * 4
	if len(left) >= limit && len(right) >= limit {
		c := make(chan []int)
		go func() {
			c <- mergeSortParallel5(left)
		}()
		go func() {
			c <- mergeSortParallel5(right)
		}()
		return merge(<- c, <- c)
	}
	return merge(mergeSort(left), mergeSort(right))
}

func timeIt(f func()) {
	s := time.Now()
	f()
	fmt.Println(time.Now().Sub(s))
}


func randomList(size int) []int {
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = rand.Intn(size)
	}
	return result
}