package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"github.com/campoy/whispering-gophers/util"
)

func main() {
	lisn, err := util.Listen()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", lisn.Addr())

	for {
		c, err := lisn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server(c)
	}
}

type Message struct {
	Addr string
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
		fmt.Printf("%+v\n", message)
	}
}
