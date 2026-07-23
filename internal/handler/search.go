package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go-snappymail/internal/config"
	"go-snappymail/internal/imap"
	"go-snappymail/internal/session"
)

// SearchHandler handles full-text search across IMAP mailboxes.
type SearchHandler struct {
	cfg *config.Config
}

func (h *SearchHandler) Results(c *echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "q is required"})
	}

	mailbox := c.QueryParam("mailbox")
	if mailbox == "" {
		mailbox = "INBOX"
	}

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SelectMailbox(mailbox); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mailbox not found"})
	}

	criteria := &imap.SearchCriteria{
		Subject: q,
		From:    q,
		Body:    q,
	}
	if c.QueryParam("unseen") == "1" {
		criteria.Unseen = true
	}

	uids, err := conn.Search(criteria)
	if err != nil {
		return err
	}

	for i, j := 0, len(uids)-1; i < j; i, j = i+1, j-1 {
		uids[i], uids[j] = uids[j], uids[i]
	}

	envelopes, err := conn.FetchEnvelopes(uids)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"mailbox":  mailbox,
		"messages": envelopes,
		"query":    q,
		"total":    len(uids),
	})
}
