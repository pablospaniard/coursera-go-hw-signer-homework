package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ExecutePipeline ...
func ExecutePipeline(in ...job) {
	wg := &sync.WaitGroup{}
	pin := make(chan interface{}, 1)
	for _, entry := range in {
		wg.Add(1)
		var pout = make(chan interface{}, 1)
		go func(wg *sync.WaitGroup, entry job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			entry(in, out)
		}(wg, entry, pin, pout)
		pin = pout
	}
	wg.Wait()
}

// func ExecutePipeline(in ...job) {
// 	c := make(chan interface{}, 100)
// 	pin := c
// 	pout := c
// 	wg := &sync.WaitGroup{}
// 	for _, entry := range in {
// 		pout = make(chan interface{}, 100)
// 		wg.Add(1)

// 		go func(wg *sync.WaitGroup, entry job, in, out chan interface{}) {
// 			defer wg.Done()
// 			defer close(out)

// 			entry(in, out)
// 		}(wg, entry, pin, pout)

// 		pin = pout
// 	}
// 	wg.Wait()
// }

// CombineResults ...
func CombineResults(in, out chan interface{}) {
	var sortedInput []string
	for input := range in {
		sortedInput = append(sortedInput, input.(string))
	}
	result := strings.Join(sortedInput, "_")
	fmt.Printf("Result is here: \n%s\n\n", result)
	out <- result
}

// SingleHash ...
func SingleHash(in, out chan interface{}) {
	println("data", in)
	for input := range in {
		var str string
		fmt.Printf("single got %d\n", input)
		str = strconv.Itoa(int(input.(int)))
		out <- DataSignerCrc32(str) + DataSignerCrc32(DataSignerMd5(str))
	}
}

// MultiHash ...
func MultiHash(in, out chan interface{}) {
	for input := range in {
		var str string
		for th := 0; th <= 5; th++ {
			str += DataSignerCrc32(strconv.Itoa(th) + input.(string))
		}
		out <- str
	}
}
