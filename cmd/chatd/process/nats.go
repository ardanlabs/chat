package process

/*
https://github.com/nats-io/go-nats

# Server
go get github.com/nats-io/gnatsd

# Run the server
gnatsd

*/

// import (
// 	"bytes"
// 	"log"
// 	"time"

// 	"github.com/ardanlabs/kit/tcp"
// 	nats "github.com/nats-io/go-nats"
// 	"github.com/pkg/errors"
// )

// // Nats subjects.
// const (
// 	natsSubject = "msg" // Handling based communication.
// )

// // natsNTMsg is sent to other Tea Servers for NTS based messaging.
// type natsNTMsg struct {
// 	id string
// 	// DATA
// }

// // String implements the fmt.Stringer interface for logging.
// func (ntMsg natsNTMsg) String() string {
// 	var b bytes.Buffer

// 	// TODO

// 	return b.String()
// }

// // ledID represents the length of the UUID based string we use for the id.
// const lenID = 36

// // natsEncode encodes the ...
// func natsEncode(id string) []byte {

// 	return data
// }

// // natsDecode decodes the ...
// func natsDecode(data []byte) (/*WHAT*/, error) {

// }

// // =============================================================================

// // NATS represents a nats system from message handling.
// type NATS struct {
// 	Config NATSConfig

// 	id   string
// 	conn *nats.Conn
// 	subs map[string]*nats.Subscription
// }

// // NATSConfig represents required configuration for the nats system.
// type NATSConfig struct {
// 	Host      string
// 	TCP   *tcp.TCP
// }

// // StartNATS initializes access to a nats system.
// func StartNATS(cfg NATSConfig) (*NATS, error) {

// 	// Set nats options for connection.
// 	opts := nats.Options{
// 		Url:            cfg.Host,
// 		AllowReconnect: true,
// 		MaxReconnect:   -1,
// 		ReconnectWait:  time.Second,
// 		Timeout:        5 * time.Second,
// 	}

// 	// Connect to the specified nats server.
// 	conn, err := opts.Connect()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "connecting to NATS")
// 	}

// 	// Construct the nats value.
// 	nts := NATS{
// 		Config: cfg,
// 		id:     uuid.NewV1().String(),
// 		conn:   conn,
// 		subs:   make(map[string]*nats.Subscription),
// 	}

// 	// Declare the event handler for handling recieved messages.
// 	f := func(msg *nats.Msg) {
// 		// Process Function
// 	}

// 	// Register the event handler for each known subject.
// 	for _, subject := range []string{natsSubject} {

// 		// Subscribe to receive messages for the specified subject.
// 		sub, err := conn.Subscribe(subject, f)
// 		if err != nil {
// 			return nil, errors.Wrapf(err, "subscribing to subject : %s", subject)
// 		}

// 		// Save the subscription with its associated subject.
// 		nts.subs[subject] = sub
// 		log.Printf("nats : subject subscribed : Subject[ %s ]\n", subject)
// 	}

// 	log.Printf("nats : service started : Host[ %s ]\n", cfg.Host)
// 	return &nts, nil
// }

// // Stop shutdowns access to the nats system.
// func (nts *NATS) Stop() {
// 	if nts == nil {
// 		log.Println("nats : WARNING : nats was not initialized")
// 		return
// 	}

// 	if nts.subs != nil {

// 		// Go through each subscription and unsubscribe.
// 		for subject, subscription := range nts.subs {
// 			if err := subscription.Unsubscribe(); err != nil {
// 				log.Printf("nats : ERROR : unsubscribe : subject[ %s ] : %v\n", subject, err)
// 				continue
// 			}

// 			log.Printf("nats : unsubscribed : subject[ %s ]\n", subject)
// 		}
// 	}

// 	log.Printf("nats : service stoped : Host[ %s ]\n", nts.Config.Host)
// }

// // Send publishes the nats  to other Tea services.
// func (nts *NATS) SendMsg(sendMsg /*WHAT*/) error {

// 	log.Printf("Nats_Process : IP[ nats ] : Outbound : Sending To NATS : %v\n", ??)

// 	return nts.conn.Publish(natsSubject, natsEncode(nts.id, ??))
// }
