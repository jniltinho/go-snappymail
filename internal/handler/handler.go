package handler

import "github.com/jniltinho/go-snappymail/internal/config"

// Handlers groups HTTP handler instances for go-snappymail.
type Handlers struct {
	Auth *AuthHandler
}

// New creates P0 handlers (auth only; mail handlers added in P1).
func New(cfg *config.Config) *Handlers {
	return &Handlers{
		Auth: &AuthHandler{cfg: cfg},
	}
}
