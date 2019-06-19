package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

func main() {
	//	get network interfaces currently available
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
	wakeOnLan("127.0.0.1:80")

}

func wakeOnLan(ip string){
	//	 port to connect to
	var port_num int
	//	split the ip string to check for any port number
	ip_addr := strings.Split(ip,":")
	//	if port is given, store it in port_num or
	//	else assign port_num a default port value (UDP 9)
	if len(ip_addr[1]) > 0{
		port_num,_ = strconv.Atoi(ip_addr[1])
	}else{
		port_num = 9
	}
	//resolve the IP address
	addr, err := net.ResolveIPAddr("ip",ip_addr[0])
	if err != nil{
		fmt.Println("Resolution error",err.Error())
		os.Exit(1)
	}
	fmt.Println(addr.IP)
	fmt.Println(port_num)
}
