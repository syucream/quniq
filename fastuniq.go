package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const (
	BufSize       = 65536
	HighWaterMark = 4
)

func process(buf *[]string, wg *sync.WaitGroup, mux *sync.Mutex, result *map[string]int) {
	defer wg.Done()

	m := map[string]int{}
	for _, key := range *buf {
		m[key] = m[key] + 1
	}

	mux.Lock()
	defer mux.Unlock()
	for k, v := range m {
		(*result)[k] = (*result)[k] + v
	}
}

func main() {
	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
	result := map[string]int{}

	s := bufio.NewScanner(os.Stdin)
	for {
		buf := make([]string, BufSize)

		// read lines until buffer is full
		isContnue := true
		for i := 0; i < cap(buf); i++ {
			isContnue = s.Scan()
			if !isContnue {
				break
			}
			buf = append(buf, s.Text())
		}

		// wait if number of goroutines is high
		if runtime.NumGoroutine() >= HighWaterMark {
			wg.Wait()
		}

		// process buffer
		wg.Add(1)
		go process(&buf, wg, mux, &result)

		if !isContnue {
			break
		}
	}
	wg.Wait()

	for k, v := range result {
		fmt.Printf("%7d %s\n", v, k)
	}
}
