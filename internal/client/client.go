package client

import (
	"net"
	"net/http"
)

type Client struct {
	localClient   *http.Client
	controlServer net.Conn
	forward       string
}

func NewClient(controlServer net.Conn, localClient *http.Client, forward string) *Client {
	client := &Client{
		controlServer: controlServer,
		localClient:   localClient,
		forward:       forward,
	}
	return client
}
