package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davidklassen/wow-pow/pkg/client"
)

var (
	addr    = flag.String("addr", "localhost:1111", "server address")
	n       = flag.Int("n", 1, "request number")
	c       = flag.Int("c", 1, "request concurrency")
	verbose = flag.Bool("v", false, "print quotes")
)

var (
	reqCount atomic.Uint64
	errCount atomic.Uint64
)

func worker(ch chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	cl := client.New(*addr)
	if err := cl.Connect(); err != nil {
		panic(err)
	}

	for range ch {
		q, err := cl.Quote()
		if err != nil {
			log.Printf("failed to fetch quote: %v", err)
			errCount.Add(1)
		}
		if *verbose {
			fmt.Println(q)
		}
		reqCount.Add(1)
	}
}

func printStats(start time.Time) {
	dur := time.Since(start)
	fmt.Println("total requests:", reqCount.Load())
	fmt.Println("error requests:", errCount.Load())
	fmt.Println("total duration:", dur)
}

func main() {
	flag.Parse()

	start := time.Now()
	ch := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(*c)
	for i := 0; i < *c; i++ {
		go worker(ch, &wg)
	}

	for i := 0; i < *n; i++ {
		ch <- struct{}{}
	}
	close(ch)
	wg.Wait()
	printStats(start)
}
