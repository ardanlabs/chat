package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ardanlabs/chat/cmd/chatd/process"
	"github.com/ardanlabs/chat/internal/platform/cache"
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tcp"
)

/*
Deal with reading partial bytes in Read call.
*/

// Configuation settings.
const (
	configKey = "CHAT"
)

func init() {
	os.Setenv("CHAT_HOST", ":6000")

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
	// Init the caching system.

	cc := cache.New()

	// =========================================================================
	// Init the socket system.

	evtFunc := func(evt, typ int, ipAddress string, format string, a ...interface{}) {
		process.Event(cc, evt, typ, ipAddress, format, a)
	}

	cfg := tcp.Config{
		NetType: "tcp4",
		Addr:    host,

		ConnHandler: process.ConnHandler{},
		ReqHandler:  process.ReqHandler{CC: cc},
		RespHandler: process.RespHandler{},

		OptEvent: tcp.OptEvent{
			Event: evtFunc,
		},
	}

	// Create a new TCP value.
	t, err := tcp.New("Sample", cfg)
	if err != nil {
		log.Printf("main : %s", err)
		return
	}

	// Start accepting client data.
	if err := t.Start(); err != nil {
		log.Printf("main : %s", err)
		return
	}
	defer t.Stop()

	log.Printf("main : Waiting for data on: %s", t.Addr())

	// =========================================================================
	// System started.

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
