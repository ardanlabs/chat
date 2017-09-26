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
Start NATS:
gnatsd

Start the Service:
CHAT_HOST=":6000" ./chatd
/*

Things TODO:
1. Add support for signaling out a specific receipient
2. Fix the Add call to cache so we don't add twice
	* Maybe add a type to the message
		* Type 1: Init
		* Type 2: Message
*/

// Configuation settings.
const configKey = "CHAT"

func init() {

	// Setup default values that can be overridden in the env.
	if _, b := os.LookupEnv("CHAT_HOST"); !b {
		os.Setenv("CHAT_HOST", ":6000")
	}
	if _, b := os.LookupEnv("CHAT_NATS_HOST"); !b {
		os.Setenv("CHAT_NATS_HOST", "nats://localhost:4222")
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
	nats := cfg.MustString("NATS_HOST")

	// =========================================================================
	// Init the caching system.

	cc := cache.New()

	// =========================================================================
	// Init the socket system.

	evtFunc := func(evt, typ int, ipAddress string, format string, a ...interface{}) {
		process.Event(cc, evt, typ, ipAddress, format, a...)
	}

	reqHandler := process.ReqHandler{
		CC: cc,
	}

	cfg := tcp.Config{
		NetType: "tcp4",
		Addr:    host,

		ConnHandler: process.ConnHandler{},
		ReqHandler:  &reqHandler,
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
	// Init NATS.

	natsCfg := process.NATSConfig{
		Host: nats,
		CC:   cc,
		TCP:  t,
	}

	nts, err := process.StartNATS(natsCfg)
	if err != nil {
		log.Printf("main : %s", err)
		return
	}
	defer nts.Stop()

	// Set our NATS access for the request handler.
	reqHandler.NATS = nts

	// =========================================================================
	// System started.

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
