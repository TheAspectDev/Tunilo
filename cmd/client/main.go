package main

import (
	"bufio"
	"bytes"
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
}

func main() {
	conn, err := net.Dial("tcp", CONTROL_SERVER_ADDRESS)
	if err != nil {
		log.Println("error while connecting to ", CONTROL_SERVER_ADDRESS)
	}
	defer conn.Close()

	err = protocol.Write(conn, protocol.Message{
		Type:      protocol.MsgReady,
		RequestID: 0,
	})

	fmt.Println(err)

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
				log.Println("error processing request")
			}

			forwardRequest(conn, request)
		}

	}
}

func forwardRequest(conn net.Conn, req *http.Request) {
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
		fmt.Fprintf(conn, "HTTP/1.1 503 Service Unavailable\r\nContent-Length: 0\r\n\r\n")
		return
	}

	defer localResp.Body.Close()

	var RequestBuffer bytes.Buffer

	if err := localResp.Write(&RequestBuffer); err != nil {
		log.Printf("Failed to serialize HTTP response: %v", err)
		return
	}
	fmt.Println(localResp.StatusCode)

	protocol.Write(conn, protocol.Message{
		Type:    protocol.MsgResponse,
		Payload: RequestBuffer.Bytes(),
	})

}
