package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ExecutePipeline ...
func ExecutePipeline(in ...job) {
	var pin = make(chan interface{}, 1)
	for _, entry := range in {
		var pout = make(chan interface{}, 1)
		go entry(pin, pout)
		pin = pout
		close(pout)
	}
}

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
		fmt.Printf("got %d\n", input)
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
