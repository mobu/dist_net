package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"regexp"
	"log"
)

//	A MAC address is 6 bytes
type MACAddress [6]byte

// A MagicPacket is constituted of 6 bytes of 0xFF followed by
// 16 groups of the destination MAC address.
type MagicPacket struct{
	header	[6]byte
	payload	[16]MACAddress
}

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
	wakeOnLan("localhost","00-50-B6-9F-BD-F6")

}

func wakeOnLan(ip string,mac string) {
	//	variable of type MagicPacket
	var packet MagicPacket
	//	variable of type MACAddress
	var macAddr MACAddress
	//	 port to connect to
	var port_num int
	//	regex statement
	//	this regex can match both 48-bit, 64-bit, and 20-octet MAC addresses
	//	Explanation of the regex below:
	//	1. ^  - this denotes that regex starts from the beginning
	//	2. () - create a capture group
	//	3. [0-9a-fA-F] - as long as anything inside the list matches
	//	4. {2} - they have to match twice
	//	5. [':'] - match the delimiter literally (exactly)
	//	6. {5} - match five consecutive patterns like above
	//	7. For the last two hex digits, we are pretty much doing the same thing except
	//	   eliminating the colon
	re_MAC := regexp.MustCompile(`^(([\da-fA-F]{2}[-:.]){5}[\da-fA-F]{2})$|^([\da-fA-F]{12})$|^(([\da-fA-F]{4}[-:.]){3}[\da-fA-F]{4})$|^(([\da-fA-F]{2}[-:.]){19}[\da-fA-F]{2})$|^(([\da-fA-F]{4}[-:.]){2}[\da-fA-F]{4})$|^(([\da-fA-F]{4}[-:.]){9}[\da-fA-F]{4})$`)
	//	if MAC address is not valid
	if !re_MAC.MatchString(mac){
		fmt.Println("MAC address" + mac + " is not valid")
	}
	//	HardwareAddr is a byte string
	hwAddr,err := net.ParseMAC(mac)
	if err != nil{
		fmt.Println("Could not parse MAC address. Please make sure it is valid.")
	}
	// copy bytes from hwAddr to macAddr bytes of MACAddress struct
	for idx := range macAddr{
		macAddr[idx] = hwAddr[idx]
	}
	//	setup the header which is 6 repetitions of 0xFF
	for idx := range packet.header{
		packet.header[idx] = 0xFF
	}
	//	setup the payload which is 16 repetitions of the MAC address
	for idx := range packet.payload{
		packet.payload[idx] = macAddr
	}
	// Fill our byte buffer with the bytes in our MagicPacket
	var buf bytes.Buffer
	if(binary.Write(&buf,binary.BigEndian,packet) != nil){
		fmt.Println("Failed writing to the buffer")
		os.Exit(1)
	}
	//	split the ip string to check for any port number
	ip_addr := strings.Split(ip, ":")
	//	if port is given, store it in port_num or
	//	else assign port_num a default port value (UDP 9)
	if len(ip_addr) > 1 {
		//	convert from string to integer
		port_num, _ = strconv.Atoi(ip_addr[1])
	} else {
		//	default UDP port number for WakeOnLan service
		port_num = 9
	}
	//	resolve the IP address
	//	here im just concatenating the ip address with the port number ip:port
	addr, err := net.ResolveUDPAddr("udp", ip_addr[0]+":"+strconv.Itoa(port_num))
	if err != nil {
		fmt.Println("Unable to get a UDP address for %s\n", addr,err.Error())
		os.Exit(1)
	}
	//	connect to the IP address
	//	 keep the local address nil
	conn,err := net.DialUDP("udp",nil,addr)
	if err != nil{
		fmt.Println("Could not connect to the destination: %s\n",addr,err.Error())
		os.Exit(1)
	}
	//	when all is done, close the connection
	defer conn.Close()
	// Write the bytes of the MagicPacket to the connection
	bytesWritten,err := conn.Write(buf.Bytes())
	if err != nil{
		fmt.Println("Unable to send packet to the destination\n")
		log.Println(err.Error())
	}else if bytesWritten != 102{
		fmt.Printf("Status: %d bytes written, %d expected\n",bytesWritten,102)
		log.Println(err.Error())
	}
}
