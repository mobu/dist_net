package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	//"strings"
)

func main() {
	interfaces, err := net.Interfaces()
	//	error handling
	if err != nil {
		fmt.Println("Error in detecting network interfaces: " + err.Error())
		os.Exit(0)
	}
	if len(interfaces) > 0 {
		fmt.Println("\nList of available network interfaces:")
		for index, i := range interfaces {
			fmt.Printf("%d.%v\n", index, i.Name)
		}
		fmt.Printf("\nSelect the interface that you want to use: (0-%d): ", len(interfaces))
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(char)
	}
}
