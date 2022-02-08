package main

import (
	"fmt"
	"math/rand"
	"time"
	// "runtime"
)

// Data -> A -> B -> C -> D
func main() {
	array := randomList(1000)
	f := func(x int) int {
		for i := 0; i < 1000_0000; i++ {
			x = x + 1
		}
		return x
	}
	timeIt("Baseline", func() {
		a := Array(array)
		a.Map(f).Map(f)
	})
	timeIt("Pipeline", func() {
		a := NewArrayPipeline(array)
		a.MapPipeline(f).MapPipeline(f).Collect()
	})
	timeIt("Parallel", func() {
		a := Array(array)
		a.MapParallel(f).MapParallel(f)
	})
}

type Array []int

func (a Array) Map(f func(int) int) Array {
	result := make(Array, len(a))
	for i := range a {
		result[i] = f(a[i])
	}
	return result
}

func (a Array) MapParallel(f func(int) int) Array {
	result := make(Array, len(a))
	mid := len(a) / 2
	
	// 2 is not important here, could be 0
	c := make(chan struct{})
	go func() {
		for i := 0; i < mid; i++ {
			result[i] = f(a[i])
		}
		c <- struct{}{}
	}()
	go func() {
		for i := mid; i < len(a); i++ {
			result[i] = f(a[i])
		}
		c <- struct{}{}
	}()
	<-c
	<-c

	return result
}

type ArrayPipeline struct {
	c chan int
	l int
}

func NewArrayPipeline(in []int) ArrayPipeline {
	out := make(chan int, 1024 * 1024 * 100)
	go func() {
		for _, x := range in {
			out <- x
		}
		close(out)
	}()
	return ArrayPipeline{
		c: out,
		l: len(in),
	}
}

func (in ArrayPipeline) MapPipeline(f func(int) int) ArrayPipeline {
	out := make(chan int, 1024 * 1024 * 100)
	go func() {
		for x := range in.c {
			out <- f(x)
		}
		close(out)
	}()
	return ArrayPipeline{
		c: out,
		l: in.l,
	}
}

func (in ArrayPipeline) Collect() []int {
	var result []int
	timeIt("allocate", func() {
		result = make([]int, in.l)
	})
	i := 0
	for x := range in.c {
		result[i] = x
		i++
	}
	return result
}

func timeIt(name string, f func()) {
	fmt.Println(name, ":")
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

func split(list []int) ([]int, []int){
	mid := len(list) / 2
	return list[0:mid], list[mid:]
}