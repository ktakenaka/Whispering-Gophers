package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"fmt"
	"flag"
	"github.com/campoy/whispering-gophers/util"
)

var (
	self string
	address = flag.String("address", "", "where to send messages")
)

type Message struct {
	Addr string
	Body string
}

func main() {
	lisn, err := util.Listen()
	if err != nil {
		log.Fatal(err)
	}
	
	self = lisn.Addr().String()
	
	log.Println("Listening on", self)
	
	flag.Parse()
	go dial(*address)

	for {
		c, err := lisn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server(c)
	}
}

func dial(addr string) {
	fmt.Println(addr)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connection established")
	}

	enc := json.NewEncoder(conn)
	stdin := bufio.NewScanner(os.Stdin)

	for stdin.Scan() {
		message := Message{Addr: self, Body: stdin.Text()}
		err := enc.Encode(message)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("OK!")
		}
	}
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
		fmt.Printf("%#v\n", message)
	}
}

