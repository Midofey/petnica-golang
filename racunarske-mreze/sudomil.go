package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "unicode"
)

func main() {
    cn, err := net.Dial("tcp", "10.100.0.98:8000")
    if err != nil {
        panic(1)
    }
    defer cn.Close()

    fmt.Printf("Uspesno povezan na SudomilChat!\nUkucaj svoje ime: ")
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    fmt.Fprintf(cn, text+"\n")

    go recieveMessage(cn)
    for {
        reader := bufio.NewReader(os.Stdin)
        text, _ := reader.ReadString('\n')
        text = strings.TrimRightFunc(text, unicode.IsSpace)
        if text != "\n" && text != "" {
            fmt.Fprintf(cn, text+"\n")
        }
    }
}

func recieveMessage(cn net.Conn) {
    reader := bufio.NewReader(cn)
    for {
        text, _ := reader.ReadString('\n')

        fmt.Printf(text)
    }
}