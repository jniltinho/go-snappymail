package handler

import "go-snappymail/internal/config"

// Handlers groups HTTP handler instances for go-snappymail.
type Handlers struct {
	Auth    *AuthHandler
	Mailbox *MailboxHandler
	Message *MessageHandler
}

// New creates application HTTP handlers.
func New(cfg *config.Config) *Handlers {
	return &Handlers{
		Auth:    &AuthHandler{cfg: cfg},
		Mailbox: &MailboxHandler{cfg: cfg},
		Message: &MessageHandler{cfg: cfg},
	}
}
