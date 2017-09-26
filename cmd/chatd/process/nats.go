package process

/*
https://github.com/nats-io/go-nats

# Server
go get github.com/nats-io/gnatsd

# Run the server
gnatsd

*/

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ardanlabs/chat/internal/msg"
	"github.com/ardanlabs/kit/tcp"
	nats "github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

// Nats subjects.
const (
	natsSubject = "msg" // Handling based communication.
)

// natsMsg is sent to other chat servers.
type natsMsg struct {
	ID  string
	MSG msg.MSG
}

// =============================================================================

// NATS represents a nats system from message handling.
type NATS struct {
	Config NATSConfig

	conn *nats.Conn
	subs map[string]*nats.Subscription
}

// NATSConfig represents required configuration for the nats system.
type NATSConfig struct {
	Host string
	ID   string
	TCP  *tcp.TCP
}

// StartNATS initializes access to a nats system.
func StartNATS(cfg NATSConfig) (*NATS, error) {

	// Set nats options for connection.
	opts := nats.Options{
		Url:            cfg.Host,
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  time.Second,
		Timeout:        5 * time.Second,
	}

	// Connect to the specified nats server.
	conn, err := opts.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "connecting to NATS")
	}

	// Construct the nats value.
	nts := NATS{
		Config: cfg,
		conn:   conn,
		subs:   make(map[string]*nats.Subscription),
	}

	// Declare the event handler for handling recieved messages.
	f := func(msg *nats.Msg) {
		// TODO: Don't process your own message, check ID.

		// Process Function
		log.Println(string(msg.Data))
	}

	// Register the event handler for each known subject.
	for _, subject := range []string{natsSubject} {

		// Subscribe to receive messages for the specified subject.
		sub, err := conn.Subscribe(subject, f)
		if err != nil {
			return nil, errors.Wrapf(err, "subscribing to subject : %s", subject)
		}

		// Save the subscription with its associated subject.
		nts.subs[subject] = sub
		log.Printf("nats : subject subscribed : Subject[ %s ]\n", subject)
	}

	log.Printf("nats : service started : Host[ %s ]\n", cfg.Host)
	return &nts, nil
}

// Stop shutdowns access to the nats system.
func (nts *NATS) Stop() {
	if nts == nil {
		log.Println("nats : WARNING : nats was not initialized")
		return
	}

	if nts.subs != nil {

		// Go through each subscription and unsubscribe.
		for subject, subscription := range nts.subs {
			if err := subscription.Unsubscribe(); err != nil {
				log.Printf("nats : ERROR : unsubscribe : subject[ %s ] : %v\n", subject, err)
				continue
			}

			log.Printf("nats : unsubscribed : subject[ %s ]\n", subject)
		}
	}

	log.Printf("nats : service stoped : Host[ %s ]\n", nts.Config.Host)
}

// SendMsg publishes the nats  to other Tea services.
func (nts *NATS) SendMsg(msg natsMsg) error {

	log.Printf("Nats_Process : IP[ nats ] : Outbound : Sending To NATS : %v\n", msg)

	data, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	return nts.conn.Publish(natsSubject, data)
}
