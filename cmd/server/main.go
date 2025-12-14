package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TheAspectDev/tunio/internal/server"
	"github.com/TheAspectDev/tunio/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	pass := flag.String("password", "12345", "Authentication password")
	controlAddr := flag.String("control", "0.0.0.0:9090", "control server address")
	publicAddr := flag.String("public", "0.0.0.0:4311", "public server address")
	flag.Parse()

	srv := server.NewServer(*publicAddr, *controlAddr, *pass)
	go srv.StartControlServer()

	http.HandleFunc("/", srv.HandleHTTP)

	go func() {
		err := http.ListenAndServe(*publicAddr, nil)

		if err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()
	lipgloss.DefaultRenderer().Output().ClearScreen()

	p := tea.NewProgram(components.SpinnerModel(*publicAddr, *controlAddr))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
