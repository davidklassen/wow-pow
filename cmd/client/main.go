package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/davidklassen/wow-pow/pkg/client"
)

var (
	addr    = flag.String("addr", "localhost:1111", "server address")
	n       = flag.Int("n", 1, "request number")
	c       = flag.Int("c", 1, "request concurrency")
	verbose = flag.Bool("v", false, "print quotes")
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
			panic(err)
		}
		if *verbose {
			fmt.Println(q)
		}
	}
}

func main() {
	flag.Parse()

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
}
