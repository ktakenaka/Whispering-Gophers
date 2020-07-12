package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	lisn, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := lisn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server(c)
	}
}

type Message struct {
	Body string
}

func server(c net.Conn) {
	defer c.Close()
	dec := json.NewDecoder(c)
	for {
		var message Message
		if err := dec.Decode(&message); err != nil {
			log.Println(err)
			return
		}
		fmt.Println(message.Body)
	}
}
