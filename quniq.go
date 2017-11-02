package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	InputBufSize        = 1 * 1024 * 1024
	OutputMapSize       = 1 * 1024 * 1024
	IntermediateMapSize = 1 * 1024
)

// print mode
type PrintMode int

const (
	PrintOnlyUnique PrintMode = iota
	PrintOnlyDuplicated
	PrintBoth
)

// memory pools
var inputPool = sync.Pool{
	New: func() interface{} {
		return make([]string, 0, InputBufSize)
	},
}

func process(buf []string, wg *sync.WaitGroup, mux *sync.Mutex, result *map[string]int, caseInsentive bool) {
	defer wg.Done()

	// pre-uniq per goroutines
	m := make(map[string]int, IntermediateMapSize)
	for _, k := range buf {
		key := k
		if caseInsentive {
			key = strings.ToLower(k)
		}

		if v, ok := m[key]; ok {
			m[key] = v + 1
		} else {
			m[key] = 1
		}
	}
	inputPool.Put(buf)

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

func printResult(result map[string]int, enableCount bool, mode PrintMode) {
	// get width of padding for count value
	padWidth := 0
	if enableCount {
		max := 0
		for _, v := range result {
			if v > max {
				max = v
			}
		}
		padWidth = len(strconv.Itoa(max))
	}

	for k, v := range result {
		if (mode == PrintOnlyUnique && v != 1) || (mode == PrintOnlyDuplicated && v == 1) {
			continue
		}

		if enableCount {
			fmt.Printf("%*d %s\n", padWidth, v, k)
		} else {
			fmt.Printf("%s\n", k)
		}
	}
}

func main() {
	// flags
	enableCount := flag.Bool("c", false, "print with count")
	onlyUnique := flag.Bool("u", false, "output only uniuqe lines")
	onlyDuplicated := flag.Bool("d", false, "output only duplicated lines")
	enableCaseInsentive := flag.Bool("i", false, "enable case insentive comparison")
	maxWorkers := flag.Int("max-workers", 1, "number of max workers")
	flag.Parse()

	// sync primitives and a shared map
	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
	result := make(map[string]int, OutputMapSize)

	s := bufio.NewScanner(os.Stdin)
	for {
		buf := inputPool.Get().([]string)
		buf = buf[:0]

		// read lines until buffer is full
		isContnue := true
		for i := 0; i < cap(buf); i++ {
			isContnue = s.Scan()
			if !isContnue {
				break
			}
			buf = append(buf, s.Text())
		}

		// process buffer
		wg.Add(1)
		go process(buf, wg, mux, &result, *enableCaseInsentive)

		// wait if number of goroutines reach max workers
		if runtime.NumGoroutine() >= *maxWorkers {
			wg.Wait()
		}

		if !isContnue {
			break
		}
	}
	wg.Wait()

	mode := PrintBoth
	if *onlyUnique && !(*onlyDuplicated) {
		mode = PrintOnlyUnique
	} else if !(*onlyUnique) && *onlyDuplicated {
		mode = PrintOnlyDuplicated
	}

	printResult(result, *enableCount, mode)
}
