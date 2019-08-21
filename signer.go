package main

import (
	"fmt"
	"sort"
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

// CombineResults ...
func CombineResults(in, out chan interface{}) {
	var sliceOfStrings []string
	for input := range in {
		sliceOfStrings = append(sliceOfStrings, input.(string))
	}
	sort.Strings(sliceOfStrings)
	result := strings.Join(sliceOfStrings, "_")
	fmt.Printf("Result is here: \n%s\n\n", result)
	out <- result
}

// CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/), объединяет отсортированный результат через _ (символ подчеркивания) в одну строку

// SingleHash ...
func SingleHash(in, out chan interface{}) {
	for input := range in {
		var str string
		str = strconv.Itoa(int(input.(int)))
		out <- DataSignerCrc32(str) + "~" + DataSignerCrc32(DataSignerMd5(str))
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
