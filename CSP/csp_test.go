package CSP

import (
	"fmt"
	"log"
	"testing"
	"unsafe"
)

func Test_singleton(t *testing.T) {
	i1 := GetSingleton("i1")
	i2 := GetSingleton("i2")
	log.Println(unsafe.Pointer(i1) == unsafe.Pointer(i2))
}

func Test_concurrency(t *testing.T) {
	job := make(chan int, 50)
	result := make(chan int, 50)
	go worker(job, result)
	for i := 1; i <= 50; i++ {
		job <- i
	}
	close(job) // 当chan被close后，for range 将取完剩余的数据后结束
	for j := 1; j <= 50; j++ {
		fmt.Println(<-result)
	}
	close(result)
}

func worker(job <-chan int, result chan<- int) {
	for n := range job {
		result <- fib(n)
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
