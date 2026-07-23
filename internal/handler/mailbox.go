package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	goimap "github.com/emersion/go-imap/v2"
	"github.com/labstack/echo/v5"
	"go-snappymail/internal/config"
	"go-snappymail/internal/imap"
	"go-snappymail/internal/session"
)

// MailboxHandler handles IMAP folder and message-list operations.
type MailboxHandler struct {
	cfg *config.Config
}

func (h *MailboxHandler) resolveCreateDelimiter(conn *imap.Client, parent, requested string) string {
	if parent == "" {
		if requested != "" {
			return requested
		}
		return "/"
	}

	folders, err := conn.ListMailboxes()
	if err == nil {
		for _, folder := range folders {
			if folder.Name == parent && folder.Delim != "" {
				return folder.Delim
			}
		}
		for _, folder := range folders {
			if folder.Delim != "" {
				return folder.Delim
			}
		}
	}

	if requested != "" {
		return requested
	}
	return "/"
}

func (h *MailboxHandler) List(c *echo.Context) error {
	mailbox := c.Param("mailbox")
	s := c.Get("imap_session").(*session.IMAPSession)

	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SelectMailbox(mailbox); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mailbox not found"})
	}

	criteria := &imap.SearchCriteria{}
	if c.QueryParam("unseen") == "1" {
		criteria.Unseen = true
	}
	if c.QueryParam("flagged") == "1" {
		criteria.Flagged = true
	}
	if q := c.QueryParam("q"); q != "" {
		criteria.Subject = q
		criteria.From = q
		criteria.Body = q
	}

	uids, err := conn.Search(criteria)
	if err != nil {
		return err
	}

	for i, j := 0, len(uids)-1; i < j; i, j = i+1, j-1 {
		uids[i], uids[j] = uids[j], uids[i]
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = h.cfg.UI.RowsPerPage
	}
	start := (page - 1) * limit
	end := start + limit
	if start > len(uids) {
		start = len(uids)
	}
	if end > len(uids) {
		end = len(uids)
	}

	fetched, err := conn.FetchEnvelopes(uids[start:end])
	if err != nil {
		return err
	}

	envMap := make(map[goimap.UID]imap.Envelope)
	for _, e := range fetched {
		envMap[e.UID] = e
	}
	envelopes := make([]imap.Envelope, 0, len(uids[start:end]))
	for _, uid := range uids[start:end] {
		if e, ok := envMap[uid]; ok {
			envelopes = append(envelopes, e)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"mailbox":  mailbox,
		"messages": envelopes,
		"page":     page,
		"limit":    limit,
		"total":    len(uids),
		"username": s.Username,
	})
}

func (h *MailboxHandler) FoldersJSON(c *echo.Context) error {
	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	folders, err := conn.ListMailboxes()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, folders)
}

func (h *MailboxHandler) CreateSubfolder(c *echo.Context) error {
	parent := c.FormValue("parent")
	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		slog.Error("IMAP connection failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create folder."})
	}
	defer conn.Close()

	delim := h.resolveCreateDelimiter(conn, parent, c.FormValue("delim"))
	fullName := name
	if parent != "" {
		fullName = parent + delim + name
	}

	if err := conn.CreateMailbox(fullName); err != nil {
		slog.Error("Failed to create mailbox", "name", fullName, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create folder."})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok", "name": fullName})
}

func (h *MailboxHandler) RenameFolder(c *echo.Context) error {
	name := c.FormValue("name")
	newname := c.FormValue("newname")
	if name == "" || newname == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and newname are required"})
	}

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		slog.Error("IMAP connection failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to rename folder."})
	}
	defer conn.Close()

	if err := conn.RenameMailbox(name, newname); err != nil {
		slog.Error("Failed to rename mailbox", "name", name, "newname", newname, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to rename folder."})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MailboxHandler) DeleteFolder(c *echo.Context) error {
	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		slog.Error("IMAP connection failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete folder."})
	}
	defer conn.Close()

	folders, err := conn.ListMailboxes()
	if err != nil {
		slog.Error("Failed to list mailboxes", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete folder."})
	}
	for _, f := range folders {
		if f.Name == name && f.IsSystem {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete system folder"})
		}
	}

	if err := conn.DeleteMailboxRecursive(name); err != nil {
		slog.Error("Failed to delete mailbox", "name", name, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete folder."})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MailboxHandler) UnreadCountJSON(c *echo.Context) error {
	name := c.Param("name")
	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	count, err := conn.UnreadCount(name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]uint32{"unseen": count})
}
