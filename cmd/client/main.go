package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/TheAspectDev/tunio/internal/client"
)

// Note: concurrency caused extra overhead and increased latency

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

	process := client.NewClient(conn, localClient, FORWARD_ADDRESS)
	process.Authenticate(pass)

	for {
		process.HandleMessage()
	}
}
