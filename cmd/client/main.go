package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TheAspectDev/tunio/internal/client"
	"github.com/TheAspectDev/tunio/internal/client/tui"
	"github.com/TheAspectDev/tunio/internal/logging"
	"github.com/TheAspectDev/tunio/internal/protocol"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func dialControlServer(controlAddr, domain string, insecure bool) (net.Conn, error) {
	if insecure {
		plainConn, err := net.Dial("tcp", controlAddr)
		if err != nil {
			return nil, fmt.Errorf("plain dial failed: %w", err)
		}
		return plainConn, nil
	} else {
		tlsConn, err := tls.Dial("tcp", controlAddr, &tls.Config{
			MinVersion: tls.VersionTLS13,
			ServerName: domain,
		})
		if err != nil {
			return nil, fmt.Errorf("tls dial failed(use -insecure if the server is not using tls): %w", err)
		}
		return tlsConn, nil
	}
}

// Note: concurrency caused extra overhead and increased latency, so ;D no concurrency
func main() {
	pass := flag.String("password", "12345", "Authentication password")
	controlAddr := flag.String("control", "127.0.0.1:9090", "control server address")
	forrwardAddr := flag.String("forward", "http://localhost:8999", "local forward address")
	noTui := flag.Bool("notui", false, "is tui used? ( false for automation/simplicity )")
	insecure := flag.Bool("insecure", false, "is the server using tls")

	flag.Parse()
	domain := strings.Split(*controlAddr, ":")[0]

	conn, err := dialControlServer(*controlAddr, domain, *insecure)

	if err != nil {
		fmt.Println(err)
		return
	}

	protocol.EnableTCPKeepalive(conn)
	session := client.NewSession(conn, localClient, *forrwardAddr)

	if *noTui {
		session.Logger = logging.StdoutLogger{}
	} else {
		session.Logger = tui.UILogger{}
	}

	defer conn.Close()

	if err := session.Authenticate(*pass); err != nil {
		session.Logger.Errorf(err, "authentication failed:")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *noTui {
		if err := session.Run(ctx); err != nil {
			session.Logger.Errorf(err, "session ended")
		}
	} else {
		go func(session client.Session, ctx context.Context) {
			if err := session.Run(ctx); err != nil {
				session.Logger.Errorf(err, "session ended")
				os.Exit(1)
			}
		}(*session, ctx)

		lipgloss.DefaultRenderer().Output().ClearScreen()

		p := tea.NewProgram(tui.ClientModel(session))

		if _, err := p.Run(); err != nil {
			session.Logger.Errorf(err, "err:")
			os.Exit(1)
		}

	}
}
