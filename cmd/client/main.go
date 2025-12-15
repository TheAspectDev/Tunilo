package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/TheAspectDev/tunio/internal/client"
	"github.com/TheAspectDev/tunio/internal/client/tui"
	"github.com/TheAspectDev/tunio/internal/logging"
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

// Note: concurrency caused extra overhead and increased latency, so ;D no concurrency
func main() {
	pass := flag.String("password", "12345", "Authentication password")
	controlAddr := flag.String("control", "127.0.0.1:9090", "control server address")
	forrwardAddr := flag.String("forward", "http://localhost:8999", "local forward address")
	noTui := flag.Bool("notui", false, "is tui used? ( false for automation/simplicity )")

	flag.Parse()

	conn, err := net.Dial("tcp", *controlAddr)
	session := client.NewSession(conn, localClient, *forrwardAddr)

	if *noTui {
		session.Logger = logging.StdoutLogger{}
	} else {
		session.Logger = tui.UILogger{}
	}

	if err != nil {
		session.Logger.Logf("error while connecting to ", *controlAddr)
		return
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
