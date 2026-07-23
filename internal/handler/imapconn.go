package handler

import (
	"time"

	"go-snappymail/internal/config"
	"go-snappymail/internal/imap"
	"go-snappymail/internal/session"
)

func imapConn(cfg *config.Config, s *session.IMAPSession) (*imap.Client, error) {
	pass, err := s.Password(cfg.Server.SecretKey)
	if err != nil {
		return nil, err
	}
	return imap.Connect(
		s.IMAPHost,
		s.IMAPPort,
		cfg.IMAP.TLS,
		cfg.IMAP.TLSServerName,
		cfg.IMAP.InsecureSkipVerify,
		time.Duration(cfg.IMAP.TimeoutSec)*time.Second,
		s.Username,
		pass,
		cfg.Server.Debug,
	)
}
