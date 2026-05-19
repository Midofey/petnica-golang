package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Println("Listening on port 8000")
	buf := make([]byte, 1024)
	for {
		n, addr, err := ln.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received %s from %s \n", string(buf[:n]), addr)
		_, err = ln.WriteToUDP([]byte(time.Now().String()), addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
