package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
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

func ParsePortRange(p string) ([]int, error) {
	var output []int
	p = strings.TrimSpace(p)
	if p == "" {
		return output, nil
	}
	portStrings := strings.Split(p, ",")

	for _, ps := range portStrings {
		nums := strings.Split(ps, "-")
		if len(nums) == 1 {
			if strings.TrimSpace(nums[0]) == "" {
				continue
			}
			n, err := strconv.Atoi(nums[0])
			if err != nil {
				return nil, errors.New("unable to process port string, cannot convert to int")
			}
			output = append(output, n)
		}

		if len(nums) > 2 {
			return nil, errors.New("unable to process port string, invalid range")
		}

		if len(nums) == 2 {
			start, err := strconv.Atoi(nums[0])
			if err != nil {
				return nil, errors.New("unable to convert start of range to integer")
			}
			stop, err := strconv.Atoi(nums[1])
			if err != nil {
				return nil, errors.New("unable to convert stop of range to integer")
			}
			if start > stop {
				return nil, errors.New("start of range must be greater than stop")
			}
			for i := start; i <= stop; i++ {
				output = append(output, i)
			}

		}

	}
	slices.Sort(output)
	return output, nil

}

//TODO: Add flag for number of workers
//TODO: Add flag for port number / port ranges

func main() {
	portFlag := flag.String("p", "", "port or range of ports to scan (i.e. 22,80,100-200)")
	flag.Usage = printUsage
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		printUsage()
		return
	}

	address := args[0]

	portsRange, err := ParsePortRange(*portFlag)
	if err != nil {
		fmt.Printf("unable to parse port flag: %v\n", err)
		return
	}

	var openPorts []int
	portsChan := make(chan int, 100)
	resultsChan := make(chan int, 100)
	defer close(portsChan)
	defer close(resultsChan)

	for i := 0; i < cap(portsChan); i++ {
		go worker(address, portsChan, resultsChan)
	}
	if len(portsRange) == 0 {
		for i := 1; i <= 1024; i++ {
			portsRange = append(portsRange, i)
		}
	}

	go func() {
		for _, p := range portsRange {
			portsChan <- p
		}
	}()

	for range len(portsRange) {
		port := <-resultsChan
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}
}
