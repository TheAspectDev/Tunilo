package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

const CONTROL_SERVER_ADDRESS = "0.0.0.0:9090"
const FORWARD_ADDRESS = "http://localhost:8999"

var localClient = &http.Client{
	Timeout: 25 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func main() {
	pass := flag.String("password", "12345", "Authentication password")
	flag.Parse()

	conn, err := net.Dial("tcp", CONTROL_SERVER_ADDRESS)
	if err != nil {
		log.Println("error while connecting to ", CONTROL_SERVER_ADDRESS)
	}
	defer conn.Close()

	var passBuffer bytes.Buffer
	writer := bufio.NewWriter(&passBuffer)
	writer.WriteString(*pass)
	fmt.Println(passBuffer.String())

	protocol.Write(conn, protocol.Message{
		Type:      protocol.MsgReady,
		Payload:   passBuffer.Bytes(),
		RequestID: 0,
	})

	for {
		msg, err := protocol.Read(conn)
		if err != nil {
			log.Println("error reading message", err)
			return
		}

		if msg.Type == protocol.MsgRequest {
			reader := bufio.NewReader(bytes.NewReader(msg.Payload))
			request, err := http.ReadRequest(reader)

			if err != nil {
				log.Println("error processing request:", err)
				continue
			}

			forwardRequest(conn, request, msg.RequestID)
		}
	}
}

func forwardRequest(conn net.Conn, req *http.Request, req_id uint64) {
	forwardData := strings.Split(FORWARD_ADDRESS, "://")

	// something.com
	req.URL.Host = forwardData[1]
	// http or https
	req.URL.Scheme = forwardData[0]

	// something.com
	req.Host = forwardData[1]

	req.RequestURI = ""

	localResp, err := localClient.Do(req)

	if err != nil {
		log.Printf("Error forwarding request to local app: %v", err)
		protocol.Write(conn, protocol.Message{
			Type:      protocol.MsgResponse,
			RequestID: req_id,
			Payload:   []byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: 0\r\n\r\n"),
		})
		return
	}

	defer localResp.Body.Close()

	var RequestBuffer bytes.Buffer

	if err := localResp.Write(&RequestBuffer); err != nil {
		log.Printf("Failed to serialize HTTP response: %v", err)
		return
	}

	protocol.Write(conn, protocol.Message{
		Type:      protocol.MsgResponse,
		RequestID: req_id,
		Payload:   RequestBuffer.Bytes(),
	})

}
