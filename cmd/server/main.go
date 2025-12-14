package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheAspectDev/tunio/internal/server"
	"github.com/TheAspectDev/tunio/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	pass := flag.String("password", "12345", "Authentication password")
	controlAddr := flag.String("control", "0.0.0.0:9090", "control server address")
	publicAddr := flag.String("public", "0.0.0.0:4311", "public server address")
	noTui := flag.Bool("notui", false, "is tui used? ( false for automation/simplicity )")

	flag.Parse()

	srv := server.NewServer(*publicAddr, *controlAddr, *pass)
	go srv.StartControlServer()

	http.HandleFunc("/", srv.HandleHTTP)

	httpServer := &http.Server{Addr: *publicAddr, Handler: nil}

	if *noTui {

		go func() {
			log.Printf("Starting serveron %s", *publicAddr)
			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("HTTP server failed: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

	} else {
		go func() {
			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Panicf("HTTP server failed: %v", err)
			}
		}()

		lipgloss.DefaultRenderer().Output().ClearScreen()

		p := tea.NewProgram(components.ServerModel(srv))

		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}

}
