package main

import (
	"fmt"
	"sync"
	"time"
)

func serve(first chan int, second chan int) {
	for {
		<-first
		var out int
		second <- out
	}
}

func client(first chan int, second chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var counter int
	max := 1000000
	for counter <= max {
		var out int
		first <- out
		<-second
		counter++
	}
}

func main() {
	start := time.Now()
	first := make(chan int)
	second := make(chan int)
	var wg sync.WaitGroup
	go serve(first, second)
	wg.Add(1)
	go client(first, second, &wg)
	wg.Wait()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed time = ", elapsed)
}
