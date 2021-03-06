package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/campoy/whispering-gophers/util"
)

var (
	self    string
	address = flag.String("address", "", "where to send messages")
)

type Message struct {
	ID string
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
	go receive()
	go dial(*address)

	for {
		c, err := lisn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server(c)
	}
}

type Peers struct {
	m  map[string]chan<- Message
	mu sync.RWMutex
}

var peers = Peers{m: make(map[string]chan<- Message)}

func (p *Peers) Add(addr string) <-chan Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.m[addr]; ok {
		return nil
	}

	ch := make(chan Message)
	p.m[addr] = ch
	return ch
}

func (p *Peers) Remove(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if ch, ok := p.m[addr]; ok {
		close(ch)
	}
	delete(p.m, addr)
}

func (p *Peers) List() []chan<- Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	lis := make([]chan<- Message, 0, len(p.m))
	for _, ch := range p.m {
		lis = append(lis, ch)
	}
	return lis
}

var seenId = struct {
	m map[string]bool
	sync.Mutext
}{m: make(map[string]bool)}

func Seen(id string) bool {
	seenId.Lock()
	defer seenId.Unlock()
	if seenId.seen[id] {
		return true
	}
	seenId.seen[id] = true
	return false
}

func receive() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		id := util.RandomID()
		message := Message{ID: id, Addr: self, Body: stdin.Text()}
		for _, ch := range peers.List() {
			select {
			case ch <- message:
			default:
				log.Println("Send failed")
			}
		}
	}
}

func dial(addr string) {
	if addr == self {
		return
	}

	ch := peers.Add(addr)
	if ch == nil {
		return
	}
	defer peers.Remove(addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connection established")
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)

	for {
		err := enc.Encode(<-ch)
		if err != nil {
			log.Fatal(err)
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
		if Seen(message.ID) {
			continue
		}
		go dial(message.Addr)
		fmt.Printf("%#v\n", message)
	}
}
