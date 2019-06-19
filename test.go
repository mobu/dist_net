package main

import (
    //"bytes"
    "fmt"
    //"io"
    "net"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "host")
        os.Exit(1)
    }

    addr, err := net.ResolveIPAddr("ip", os.Args[1])
    if err != nil {
        fmt.Println("Resolution error", err.Error())
        os.Exit(1)
	}
	fmt.Println(addr)
}