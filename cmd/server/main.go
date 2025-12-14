package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/TheAspectDev/tunio/internal/server"
)

func main() {
	pass := flag.String("password", "12345", "Authentication password")
	controlAddr := flag.String("control", "0.0.0.0:9090", "control server address")
	publicAddr := flag.String("public", "0.0.0.0:4311", "public server address")
	flag.Parse()

	srv := server.NewServer(*publicAddr, *controlAddr, *pass)
	go srv.StartControlServer()

	http.HandleFunc("/", srv.HandleHTTP)
	err := http.ListenAndServe(*publicAddr, nil)

	if err != nil {
		log.Fatal("HTTP srvr failed:", err)
	}
}
