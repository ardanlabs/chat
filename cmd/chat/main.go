package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ardanlabs/chat/internal/msg"
	"github.com/ardanlabs/kit/cfg"
)

/*
Start the Client:
CHAT_HOST=":6000" ./chat
*/

// Configuation settings.
const configKey = "CHAT"

func init() {

	// Setup default values that can be overridden in the env.
	if _, b := os.LookupEnv("CHAT_HOST"); !b {
		os.Setenv("CHAT_HOST", ":6000")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
}

func main() {

	// =========================================================================
	// Init the configuration system.

	if err := cfg.Init(cfg.EnvProvider{Namespace: configKey}); err != nil {
		log.Println("Error initalizing configuration system", err)
		os.Exit(1)
	}

	log.Println("Configuration\n", cfg.Log())

	// Get configuration.
	host := cfg.MustString("HOST")

	// =========================================================================
	// Connect and get going.

	// Let's connect back and send a TCP package
	conn, err := net.Dial("tcp4", host)
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
