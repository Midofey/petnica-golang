package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleClient(conn net.Conn) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(message))
	conn.Write([]byte(message))
}

func main() {

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening on port 8000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleClient(conn)
	}
}
