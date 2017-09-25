package main

import (
	"bufio"
	"log"
	"net"

	"github.com/ardanlabs/chat/internal/msg"
)

func main() {

	// Let's connect back and send a TCP package
	conn, err := net.Dial("tcp4", ":6000")
	if err != nil {
		log.Println("dial", err)
	}

	m := msg.MSG{
		Name: "0123456789",
		Data: "Hello There",
	}

	data := msg.Encode(m)

	if _, err := conn.Write(data); err != nil {
		log.Println("write", err)
	}

	bufReader := bufio.NewReader(conn)
	response, err := bufReader.ReadString('\n')
	log.Println(response)
}
