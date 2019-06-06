package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	query := "Our Query"
	respond := make(chan string, 1)

	go googleIt(respond, query)

	select {
	case queryResp := <-respond:
		fmt.Printf("Sent query:\t\t %s\n", query)
		fmt.Printf("Got Response:\t\t %s\n", queryResp)

	case <-time.After(5 * time.Second):
		fmt.Printf("A timeout occurred for query:\t\t %s\n", query)
	}
}

func googleIt(respond chan<- string, query string) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	respond <- "A Google Response"
}
