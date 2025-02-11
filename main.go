package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"terminal-resume.jayash.space/models"
	templates "terminal-resume.jayash.space/templates/simple"
)

	const (
		host = "0.0.0.0" // Listen on all interfaces
		port = "2222"   // Default SSH port
	) 






	//go:embed data.json
	var jsonContent []byte

	func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

		var jsonData models.JsonData
		err := json.Unmarshal(jsonContent, &jsonData)
		var username string

		if err != nil {
			fmt.Println("error unmarshaling JSON:", err)
			os.Exit(1)
		}

		// err = extractUsername(s, &username)

		// if err != nil {
		// 	log.Warn(err)
		// }

		log.Info("Username is " + username)


		m := templates.SimpleModel{
			Sess:    s,
			Content: jsonData,
		}
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}



	func main() {

		s, err := wish.NewServer(
			wish.WithAddress(net.JoinHostPort(host, port)),
			wish.WithHostKeyPath(".ssh/id_ed25519"),
			// Allocate a pty.
			// This creates a pseudoconsole on windows, compatibility is limited in
			// that case, see the open issues for more details.
			ssh.AllocatePty(),
			wish.WithMiddleware(
				// run our Bubble Tea handler
				bubbletea.Middleware(teaHandler),
				// ensure the user has requested a tty
				activeterm.Middleware(),
				logging.Middleware(),

			),
			ssh.HostKeyFile(".ssh/id_ed25519"),
		)

		if err != nil {
			log.Error("Could not start server", "error", err)
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Info("Starting SSH server", "host", host, "port", port)

		go func() {
			if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				log.Error("Could not start server", "error", err)
				done <- nil
			}
		}()

		<-done
		log.Info("Stopping SSH server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer func() { cancel() }()
		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not stop server", "error", err)
		}

	}
