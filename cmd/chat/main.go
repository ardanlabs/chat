package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ardanlabs/chat/internal/msg"
)

/*
// Accept keyboard input.
	reader := bufio.NewReader(os.Stdin)

	// Process keyboard input.
	for {
		fmt.Print("\n#> ")
		cmd, _ := reader.ReadString('\n')
	}
*/

func main() {

	// Let's connect back and send a TCP package
	conn, err := net.Dial("tcp4", ":6000")
	if err != nil {
		log.Println("dial", err)
	}

	// Accept keyboard input.
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nName:> ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	// Show online.
	mSend := msg.MSG{
		Sender:    name,
		Recipient: "",
		Data:      fmt.Sprintf("%s is online", name),
	}
	data := msg.Encode(mSend)
	if _, err := conn.Write(data); err != nil {
		log.Println("write", err)
	}

	// Receiving goroutine.
	go func() {
		for {
			data, _, err := msg.Read(conn)
			if err != nil {
				log.Println("read", err)
				return
			}

			mRecv := msg.Decode(data)
			log.Println(mRecv)
			fmt.Printf("\n%s#> ", name)
		}
	}()

	// Process keyboard input.
	for {
		fmt.Printf("\n%s#> ", name)
		message, _ := reader.ReadString('\n')

		mSend := msg.MSG{
			Sender:    name,
			Recipient: "",
			Data:      message,
		}

		data := msg.Encode(mSend)

		if _, err := conn.Write(data); err != nil {
			log.Println("write", err)
		}
	}
}
