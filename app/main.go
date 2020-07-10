package main

import (
	"os"
	"bufio"
	"log"
	"encoding/json"
)

func main() {
	enc := json.NewEncoder(os.Stdout)
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		log.Println("accepting input")
		message := Message{stdin.Text()}
		err := enc.Encode(message)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("OK!")
		}
	}
}

type Message struct {
	Body string
}
