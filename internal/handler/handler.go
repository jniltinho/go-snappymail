package handler

import "go-snappymail/internal/config"

// Handlers groups HTTP handler instances for go-snappymail.
type Handlers struct {
	Auth    *AuthHandler
	Mailbox *MailboxHandler
	Message *MessageHandler
	Compose *ComposeHandler
	Search  *SearchHandler
	UI      *UIHandler
}

// New creates application HTTP handlers.
func New(cfg *config.Config) *Handlers {
	return &Handlers{
		Auth:    &AuthHandler{cfg: cfg},
		Mailbox: &MailboxHandler{cfg: cfg},
		Message: &MessageHandler{cfg: cfg},
		Compose: &ComposeHandler{cfg: cfg},
		Search:  &SearchHandler{cfg: cfg},
		UI:      &UIHandler{cfg: cfg},
	}
}
