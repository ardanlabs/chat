package process

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/ardanlabs/chat/internal/msg"
	"github.com/ardanlabs/chat/internal/platform/cache"
	"github.com/ardanlabs/kit/tcp"
)

// Event writes tcp events.
func Event(cc *cache.Cache, evt, typ int, ipAddress string, format string, a ...interface{}) {
	log.Printf("****> EVENT : IP[ %s ] : EVT[%s] TYP[%s] : %s", ipAddress, evtTypes[evt], typTypes[typ], fmt.Sprintf(format, a...))

	if typ == tcp.TypTrigger {
		switch evt {
		case tcp.EvtDrop:
			client, err := cc.GetAddress(ipAddress)
			if err != nil {
				log.Printf("****> EVENT : IP[ %s ] : ERROR : alread removed from cache.", ipAddress)
				return
			}

			if err := cc.Remove(ipAddress); err != nil {
				log.Printf("****> EVENT : IP[ %s ] : ERROR : removing from cache : %s", ipAddress, err)
				return
			}
			log.Printf("****> EVENT : IP[ %s ] : removed [ %s ] from cache.", ipAddress, client.ID)
		}
	}
}

// Process handles all the communication logic.
func Process(cc *cache.Cache, nats *NATS, r *tcp.Request) {
	ipAddress := r.TCPAddr.String()

	// Decode the message bytes into a msg.MSG.
	m := msg.Decode(r.Data)

	log.Printf("Socket_Process : IP[ %s ] : Inbound : %v\n", ipAddress, m)

	// Make sure this client is in the cache.
	// TODO: This will fail if the Sender is already in the cache.
	// TODO: We need a type in the message for the first message.
	// TODO: Then we can check if this is here already.
	cc.Add(m.Sender, r.TCPAddr)

	// Send the message to NATS for processing.
	if err := nats.SendMsg(m); err != nil {
		log.Printf("Socket_Process : IP[ %s ] : ERROR : %s\n", ipAddress, err)
	}
}

// =============================================================================

var evtTypes = []string{
	"unknown",
	"Accept",
	"Join",
	"Read",
	"Remove",
	"Drop",
	"Groom",
}

// Set of event sub types.
var typTypes = []string{
	"unknown",
	"Error",
	"Info",
	"Trigger",
}

// =============================================================================

// ConnHandler is required to process data.
type ConnHandler struct{}

// Bind is called to init a reader and writer.
func (ConnHandler) Bind(conn net.Conn) (io.Reader, io.Writer) {
	return conn, conn
}

// ReqHandler is required to process client messages.
type ReqHandler struct {
	CC   *cache.Cache
	NATS *NATS
}

// Read implements the tcp.ReqHandler interface. It is provided a request
// value to populate and a io.Reader that was created in the Bind above.
func (*ReqHandler) Read(ipAddress string, reader io.Reader) ([]byte, int, error) {

	// Block on the network for our message.
	data, n, err := msg.Read(reader)
	if err != nil {
		log.Printf("read : IP[ %s ] : %s", ipAddress, err)
		return nil, 0, err
	}

	log.Printf("read : IP[ %s ] : Length[%d]", ipAddress, len(data))
	return data, n, nil
}

// Process is used to handle the processing of the message. This method
// is called on a routine from a pool of routines.
func (req *ReqHandler) Process(r *tcp.Request) {
	Process(req.CC, req.NATS, r)
}

// RespHandler is required to send messages.
type RespHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (RespHandler) Write(r *tcp.Response, writer io.Writer) error {
	log.Printf("write : IP[ %s ] : Length[ %d ]\n", r.TCPAddr.IP.String(), len(r.Data))

	if _, err := writer.Write(r.Data); err != nil {
		return err
	}
	return nil
}
