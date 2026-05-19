package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func recieveMessage(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(message)
	}
}

func main() {
	conn, err := net.Dial("tcp", "10.100.0.84:8000")
	if err != nil {
		fmt.Println(err)
	}

	for {
		go recieveMessage(conn)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text+"\n")
	}
}
