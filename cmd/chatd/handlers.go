package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"

	"github.com/ardanlabs/kit/tcp"
)

// connHandler is required to process data.
type connHandler struct{}

// Bind is called to init a reader and writer.
func (connHandler) Bind(conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

// reqHandler is required to process client messages.
type reqHandler struct{}

// Read implements the tcp.ReqHandler interface. It is provided a request
// value to populate and a io.Reader that was created in the Bind above.
func (reqHandler) Read(ipAddress string, reader io.Reader) ([]byte, int, error) {
	bufReader := reader.(*bufio.Reader)

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		log.Printf("read : IP[ %s ] : %s", ipAddress, err)
		return nil, 0, err
	}

	log.Printf("read : IP[ %s ] : Length[%d]", ipAddress, len(line))
	return []byte(line), len(line), nil
}

// Process is used to handle the processing of the message. This method
// is called on a routine from a pool of routines.
func (reqHandler) Process(r *tcp.Request) {
	log.Printf("read : IP[ %s ] : %s\n", r.TCPAddr.IP.String(), string(r.Data))

	resp := tcp.Response{
		TCPAddr: r.TCPAddr,
		Data:    []byte("GOT IT\n"),
		Length:  7,
	}

	r.TCP.Send(context.TODO(), &resp)
}

type respHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (respHandler) Write(r *tcp.Response, writer io.Writer) error {
	bufWriter := writer.(*bufio.Writer)
	if _, err := bufWriter.WriteString(string(r.Data)); err != nil {
		return err
	}
	return bufWriter.Flush()
}
