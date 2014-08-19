package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
)

var origin = "http://localhost/"
var url = "ws://localhost:3000/search?q=#艦これ"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		var msg = make([]byte, 2048)
		_, err = ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Receive: %s\n", msg)
	}
}
