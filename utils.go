package main

import (
	"errors"

	"github.com/charmbracelet/ssh"
)

func extractUsername(s ssh.Session, user *string) error {
	if len(s.Command()) == 0 || s.Command()[0] == "" {
		return errors.New("no username provided")
	}
	*user = s.Command()[0]
	return nil
}