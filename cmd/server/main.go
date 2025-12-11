package main

import (
	"log"
	"net/http"

	"github.com/TheAspectDev/tunio/internal/server"
)

const CONTROL_SERVER_ADDRESS = "0.0.0.0:9090"
const PUBLIC_SERVER_ADDRESS = "0.0.0.0:4311"

func main() {

	srv := server.NewServer(PUBLIC_SERVER_ADDRESS, CONTROL_SERVER_ADDRESS)

	go srv.StartControlServer()

	http.HandleFunc("/", srv.HandleHTTP)
	err := http.ListenAndServe(PUBLIC_SERVER_ADDRESS, nil)

	if err != nil {
		log.Fatal("HTTP srvr failed:", err)
	}

}
