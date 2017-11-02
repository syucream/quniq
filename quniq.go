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
	MegaNum             = 1024 * 1024
	OutputMapSize       = 1 * 1024 * 1024
	IntermediateMapSize = 1 * 1024
)

type PrintMode int

const (
	PrintOnlyUnique PrintMode = iota
	PrintOnlyDuplicated
	PrintBoth
)

type MovableBuffer struct {
	pool sync.Pool
	buf  []string
}

func NewMovableBuffer(p sync.Pool, b []string) *MovableBuffer {
	return &MovableBuffer{
		pool: p,
		buf:  b,
	}
}

func (mbuf *MovableBuffer) move() {
	mbuf.pool.Put(mbuf.buf)
}

func process(mbuf *MovableBuffer, wg *sync.WaitGroup, mux *sync.Mutex, result *map[string]int, caseInsentive bool) {
	defer wg.Done()

	// 1st level uniq per goroutines
	m := make(map[string]int, IntermediateMapSize)
	for _, k := range mbuf.buf {
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
	mbuf.move()

	// 2nd level uniq with getting mutex
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
	inputBufferWeight := flag.Int("inbuf-weight", 1, "number of input buffer items(used specified value * 1024 * 1024)")
	maxWorkers := flag.Int("max-workers", 1, "number of max workers")
	flag.Parse()

	// sync primitives and a shared map
	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
	result := make(map[string]int, OutputMapSize)

	// input buffer pool
	bufferPool := sync.Pool{
		New: func() interface{} {
			return make([]string, 0, *inputBufferWeight*MegaNum)
		},
	}

	s := bufio.NewScanner(os.Stdin)
	for {
		// get buffer from pool and reset it
		buf := bufferPool.Get().([]string)
		buf = buf[:0]

		// read lines until buffer is full or input reaches EOF
		isContnue := true
		for i := 0; i < cap(buf); i++ {
			isContnue = s.Scan()
			if !isContnue {
				break
			}
			buf = append(buf, s.Text())
		}

		// move buffer to other goroutine and process
		mbuf := NewMovableBuffer(bufferPool, buf)
		wg.Add(1)
		go process(mbuf, wg, mux, &result, *enableCaseInsentive)

		// wait if number of goroutines reach max workers for resource limitation
		if runtime.NumGoroutine() >= *maxWorkers {
			wg.Wait()
		}

		if !isContnue {
			break
		}
	}
	wg.Wait()

	// prepare print mode
	mode := PrintBoth
	if *onlyUnique && !(*onlyDuplicated) {
		mode = PrintOnlyUnique
	} else if !(*onlyUnique) && *onlyDuplicated {
		mode = PrintOnlyDuplicated
	}

	printResult(result, *enableCount, mode)
}
