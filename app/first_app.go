package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
)

var address = flag.String("address", "", "where to send messages")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *address)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connection established")
	}

	enc := json.NewEncoder(conn)
	stdin := bufio.NewScanner(os.Stdin)

	for stdin.Scan() {
		log.Println("accepting input")
		message := Message{stdin.Text()}
		err := enc.Encode(message)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("OK!")
		}
	}
}

type Message struct {
	Body string
}
