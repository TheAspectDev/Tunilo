package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheAspectDev/tunio/internal/server"
	"github.com/TheAspectDev/tunio/internal/server/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	pass := flag.String(
		"password",
		"12345",
		"Password used for client authentication",
	)

	controlAddress := flag.String(
		"control",
		"0.0.0.0:9090",
		"Address the control server listens on (host:port)",
	)

	publicAddress := flag.String(
		"public",
		"0.0.0.0:4311",
		"Address the public server listens on (host:port)",
	)

	noTui := flag.Bool(
		"notui",
		false,
		"Disable the interactive TUI (useful for automation or headless environments)",
	)

	insecure := flag.Bool(
		"insecure",
		false,
		"Disable TLS and allow insecure connections",
	)

	cert := flag.String(
		"cert",
		"",
		"Path to the TLS certificate file (required when TLS is enabled)",
	)

	key := flag.String(
		"key",
		"",
		"Path to the TLS private key file (required when TLS is enabled)",
	)

	flag.Parse()

	srvBuilder := server.NewServerBuilder().
		SetAddress(*publicAddress).
		SetControlAddress(*controlAddress).SetPassword(*pass)

	if !*insecure {
		srvBuilder = srvBuilder.SetTLS(server.TLSConfig{
			Cert: *cert,
			Key:  *key,
		})
	}

	srv, err := srvBuilder.Build()

	httpServer := &http.Server{Addr: *publicAddress, Handler: nil}

	if !*insecure {
		httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", srv.HandleHTTP)

	go srv.StartControlServer()

	if *noTui {
		log.Printf("Starting serveron %s", *publicAddress)
		go srv.StartPublicServer(httpServer)

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

	} else {
		go srv.StartPublicServer(httpServer)

		lipgloss.DefaultRenderer().Output().ClearScreen()

		p := tea.NewProgram(tui.ServerModel(srv))

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
