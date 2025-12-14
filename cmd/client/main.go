package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/TheAspectDev/tunio/internal/client"
)

var localClient = &http.Client{
	Timeout: 25 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

// Note: concurrency caused extra overhead and increased latency, so ;D no concurrency
func main() {
	pass := flag.String("password", "12345", "Authentication password")
	controlAddr := flag.String("control", "127.0.0.1:9090", "control server address")
	forrwardAddr := flag.String("forward", "http://localhost:8999", "local forward address")
	flag.Parse()

	conn, err := net.Dial("tcp", *controlAddr)
	if err != nil {
		log.Println("error while connecting to ", *controlAddr)
	}
	defer conn.Close()

	session := client.NewSession(conn, localClient, *forrwardAddr)

	if err := session.Authenticate(*pass); err != nil {
		log.Fatal("authentication failed:", err)
	}

	// NOTE: not used yet, built for TUI: ctrl+c-quit support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := session.Run(ctx); err != nil {
		log.Println("session ended:", err)
	}
}
