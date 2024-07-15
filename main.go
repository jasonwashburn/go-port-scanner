package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
)

func worker(address string, ports chan int, results chan int) {
	for p := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, p))
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func printUsage() {
	fmt.Printf("Usage: %s [OPTIONS] <address>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		printUsage()
		return
	}

	address := args[0]

	var openPorts []int
	ports := make(chan int, 100)
	results := make(chan int, 100)
	defer close(ports)
	defer close(results)

	for i := 0; i < cap(ports); i++ {
		go worker(address, ports, results)
	}

	go func() {
		for i := 1; 1 <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i <= 1024; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}
}
