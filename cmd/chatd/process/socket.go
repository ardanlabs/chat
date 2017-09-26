package process

import (
	"context"
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
func Process(id string, cc *cache.Cache, nats *NATS, r *tcp.Request) {
	m := msg.Decode(r.Data)

	log.Printf("process : IP[ %s ] : %v\n", r.TCPAddr.IP.String(), m)

	// TODO: Handle errors for multiple adds by the same ID.
	// TODO: Add an auth message to handle this better.
	cc.Add(m.Sender, r.TCPAddr)

	d := msg.Encode(m)

	for _, client := range cc.Get(m.Sender) {
		log.Printf("process : IP[ %s ] : Client[ %s ]\n", r.TCPAddr.IP.String(), client.ID)

		resp := tcp.Response{
			TCPAddr: client.TCPAddr,
			Data:    d,
			Length:  len(d),
		}

		r.TCP.Send(context.TODO(), &resp)

		// TODO: Only do this if the client is not in our cache.
		nm := natsMsg{
			ID:  id,
			MSG: m,
		}
		if err := nats.SendMsg(nm); err != nil {
			log.Printf("process : IP[ NATS ] : ERROR : Client[ %s ]\n", client.ID)
		}
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
	ID   string
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
	Process(req.ID, req.CC, req.NATS, r)
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
