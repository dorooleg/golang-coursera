package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func Crc32(data string) chan string {
	out := make(chan string)
	go func(data string) {
		out <- DataSignerCrc32(data)
		close(out)
	}(data)
	return out
}

var mtx sync.Mutex = sync.Mutex{}

func Md5(data string) string {
	mtx.Lock()
	result := DataSignerMd5(data)
	mtx.Unlock()
	return result
}

func SingleHash(in, out chan interface{}) {
	left, right := make([]chan string, 0), make([]chan string, 0)
	for entry := range in {
		data := strconv.Itoa(entry.(int))
		left, right = append(left, Crc32(data)), append(right, Crc32(Md5(data)))
	}
	for idx := range left {
		out <- <-left[idx] + "~" + <-right[idx]
	}
}

func MultiHash(in, out chan interface{}) {
	results := make([]chan string, 0)
	for entry := range in {
		var data string = entry.(string)
		for th := 0; th < 6; th++ {
			results = append(results, Crc32(strconv.Itoa(th)+data))
		}
	}
	for i := 0; i < len(results); i += 6 {
		var result string
		for th := 0; th < 6; th++ {
			result += <-results[i+th]
		}
		out <- result
	}
}

func CombineResults(in, out chan interface{}) {
	results := make([]string, 0)
	for entry := range in {
		results = append(results, entry.(string))
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}

func ExecutePipeline(jobs ...job) {
	var in chan interface{}
	wg := &sync.WaitGroup{}
	for _, task := range jobs {
		wg.Add(1)
		out := make(chan interface{}, 1)
		go func(wg *sync.WaitGroup, task job, in, out chan interface{}) {
			defer wg.Done()
			task(in, out)
			close(out)
		}(wg, task, in, out)
		in = out
	}
	wg.Wait()
}
