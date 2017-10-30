package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const (
	InputBufSize        = 4 * 1024 * 1024
	OutputMapSize       = 1 * 1024 * 1024
	IntermediateMapSize = 1 * 1024
)

func process(buf []string, wg *sync.WaitGroup, mux *sync.Mutex, result *map[string]int) {
	defer wg.Done()

	// pre-uniq per goroutines
	m := make(map[string]int, IntermediateMapSize)
	for _, key := range buf {
		if v, ok := m[key]; ok {
			m[key] = v + 1
		} else {
			m[key] = 1
		}
	}

	// merge uniq values with getting mutex
	mux.Lock()
	defer mux.Unlock()
	for k, v := range m {
		if rv, ok := (*result)[k]; ok {
			(*result)[k] = rv + v
		} else {
			(*result)[k] = v
		}
	}
}

func main() {
	// flags
	enableCount := flag.Bool("c", false, "print with count")
	maxWorkers := flag.Int("workers", 1, "number of max workers")
	flag.Parse()

	// sync primitives and a shared map
	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
	result := make(map[string]int, OutputMapSize)

	s := bufio.NewScanner(os.Stdin)
	buf := make([]string, 0, InputBufSize)
	for {
		cursor := buf[:0]

		// read lines until buffer is full
		isContnue := true
		for i := 0; i < cap(cursor); i++ {
			isContnue = s.Scan()
			if !isContnue {
				break
			}
			cursor = append(cursor, s.Text())
		}

		// copy and process buffer
		// NOTE it has allocation and copy cost ...
		wg.Add(1)
		processBuffer := make([]string, len(cursor))
		copy(processBuffer, cursor)
		go process(processBuffer, wg, mux, &result)

		// wait if number of goroutines reach max workers
		if runtime.NumGoroutine() >= *maxWorkers {
			wg.Wait()
		}

		if !isContnue {
			break
		}
	}
	wg.Wait()

	if *enableCount {
		// similar to uniq -c
		for k, v := range result {
			fmt.Printf("%7d %s\n", v, k)
		}
	} else {
		for k, _ := range result {
			fmt.Printf("%s\n", k)
		}
	}
}
