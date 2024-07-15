package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err == nil {
		fmt.Println("Connection successful")
	}
	defer conn.Close()
}
