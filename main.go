package main

import (
	"encoding/json"
	"flag" //Package flag implements command-line flag parsing
	"fmt"  //formatting and printing
	"math/rand"
	"net"     //net provides a portable interface for network I/O
	"strconv" //implements conversions to and from string representations of basic data types
	"strings"
	"time"
)

/* Information about node */
type NodeInfo struct {
	NodeId     int    `json:"nodeId"`
	NodeIpAddr string `json:"nodeIpAddr"`
	Port       string `json:"port"`
}

// A standard format for adding node to cluster
type AddToClusterMessage struct {
	Source  NodeInfo `json:"source"`
	Dest    NodeInfo `json:"dest"`
	Message string   `json:"message"`
}

/* Just for pretty printing the node info */
/* using go's method implementation to attach String method to the
NodeInfo struct */
func (node NodeInfo) String() string {
	return "NodeInfo:{ nodeId: " + strconv.Itoa(node.NodeId) + ", nodeIpAddr: " + node.NodeIpAddr + ", port:" + node.Port + " }"
}

/* Just for pretty printing Request/Response info */
// here we are overloading the String() function
func (req AddToClusterMessage) String() string {
	return "AddToClusterMessage:{\n source: " + req.Source.String() + ",\n dest: " + req.Dest.String() + ",\n message: " + req.Message + " }"
}

func main() {
	/* get command-line arguments */
	makeMasterOnError := flag.Bool("makeMasterOnError", false, "make this node master if unable to connect to the cluster ip provided.")
	clusterip := flag.String("clusterip", "127.0.0.1:8001", "ip address of any node to connect")
	myport := flag.String("myport", "8001", "ip address to run this node on. default is 8001.")
	flag.Parse()

	/* Generate id */
	/* here we are using the rand.Seed() method to generate random seed
	for the RNG */
	// rand.Seed() expects a 64 bit integer value and here we are doing so
	// by passing the value of current time and converting the time into
	// 64 bit integer using UnixNano() function
	rand.Seed(time.Now().UTC().UnixNano())
	myid := rand.Intn(99999999)

	// InterfaceAddrs returns a list of the system's network interface addresses. 
	addr, err := net.InterfaceAddrs()
	// if there's an error, show error message and exit program
	if err != nil {
		os.Stderr.WriteString("Could not retrieve interface addresses.\nError: " + err.Error() + "\n")
		os.Exit(1)
	}
	// go through all the network interfaces found
	// range cycles through the whole array/slice of InterfaceAddrs()
	for _, a := range addr {
		// short one-liner to retrieve the value for IP address
		// the ok keyword is a boolean in Go and it is true if the function
		// returns true
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// IP.To4 converts the IPv4 address to a 4-byte representation
			if ipnet.IP.To4() != nil {
				myIp := ipnet.IP.String()
			}
		}
	}
	/* create a NodeInfo struct using the values provided
	 *myport is deferencing the pointer myport to retrieve the actual value */
	me := NodeInfo{NodeId: myid, NodeIpAddr: myIp[0].String(), Port: *myport}

	/* Still figuring out why the NodeId is -1.
	getting the IP address and Port number by using the split func on clusterip */
	dest := NodeInfo{NodeId: -1, NodeIpAddr: strings.Split(*clusterip, ":")[0], Port: strings.Split(*clusterip, ":")[1]}
	fmt.Println("My details:", me.String())

	/* Try to connect to the cluster, and send request to cluster if able to connect */
	ableToConnect := connectToCluster(me, dest)

	/* Check if you are able to connect to a cluster and also check
	if makeMasterOnError is checked or not. If no cluster is found and
	makeMasterOnError is checked, then configure this node as the master */
	if ableToConnect || (!ableToConnect && *makeMasterOnError) {
		if *makeMasterOnError {
			fmt.Println("Will start this node as master.")
		}
		listenOnPort(me)
	} else {
		fmt.Println("Quitting system. Set makeMasterOnError flag to make the node master.", myid)
	}
}

/*
 * This is a useful utility to format the json packet to send requests
 * This tiny block is sort of important else you will end up sending blank messages.
 */
func getAddToClusterMessage(source NodeInfo, dest NodeInfo, message string) AddToClusterMessage {
	return AddToClusterMessage{
		Source: NodeInfo{
			NodeId:     source.NodeId,
			NodeIpAddr: source.NodeIpAddr,
			Port:       source.Port,
		},
		Dest: NodeInfo{
			NodeId:     dest.NodeId,
			NodeIpAddr: dest.NodeIpAddr,
			Port:       dest.Port,
		},
		Message: message,
	}
}

func connectToCluster(me NodeInfo, dest NodeInfo) bool {
	/* connect to this socket details provided */
	connOut, err := net.DialTimeout("tcp", dest.NodeIpAddr+":"+dest.Port, time.Duration(10)*time.Second)

	// if an error occurred
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Println("Couldn't connect to cluster.", me.NodeId)
			return false
		}
	} else {
		fmt.Println("Connected to cluster. Sending message to node.")
		text := "Hi nody.. please add me to the cluster.."
		requestMessage := getAddToClusterMessage(me, dest, text)
		json.NewEncoder(connOut).Encode(&requestMessage)

		decoder := json.NewDecoder(connOut)
		var responseMessage AddToClusterMessage
		decoder.Decode(&responseMessage)
		fmt.Println("Got response:\n" + responseMessage.String())

		return true
	}
	return false
}

func listenOnPort(me NodeInfo) {
	/* Listen for incoming messages */
	ln, _ := net.Listen("tcp", fmt.Sprint(":"+me.Port))
	/* accept connection on port */
	/* not sure if looping infinetely on ln.Accept() is good idea */
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				fmt.Println("Error received while listening.", me.NodeId)
			}
		} else {
			var requestMessage AddToClusterMessage
			json.NewDecoder(connIn).Decode(&requestMessage)
			fmt.Println("Got request:\n" + requestMessage.String())

			text := "OK."
			responseMessage := getAddToClusterMessage(me, requestMessage.Source, text)
			json.NewEncoder(connIn).Encode(&responseMessage)
			connIn.Close()
		}
	}
}
