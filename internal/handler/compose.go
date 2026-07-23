package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	goimap "github.com/emersion/go-imap/v2"
	"github.com/labstack/echo/v5"
	"go-snappymail/internal/config"
	"go-snappymail/internal/imap"
	"go-snappymail/internal/session"
	smtppkg "go-snappymail/internal/smtp"
)

// ComposeHandler handles email composition and sending.
type ComposeHandler struct {
	cfg *config.Config
}

func (h *ComposeHandler) Send(c *echo.Context) error {
	s := c.Get("imap_session").(*session.IMAPSession)
	pass, err := s.Password(h.cfg.Server.SecretKey)
	if err != nil {
		return err
	}

	fromEmail := c.FormValue("from_email")
	if fromEmail == "" {
		fromEmail = s.Username
	}
	msg := &smtppkg.Message{
		From:        fromEmail,
		DisplayName: c.FormValue("from_name"),
		To:          splitAddrs(c.FormValue("to")),
		Cc:          splitAddrs(c.FormValue("cc")),
		Bcc:         splitAddrs(c.FormValue("bcc")),
		Subject:     c.FormValue("subject"),
		TextHTML:    c.FormValue("body_html"),
		TextPlain:   c.FormValue("body_plain"),
	}
	if replyTo := strings.TrimSpace(c.FormValue("reply_to")); replyTo != "" {
		msg.ReplyTo = replyTo
	}

	if err := h.attachFormFiles(c, msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if len(msg.To) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "to is required"})
	}

	smtpCfg := smtppkg.Config{
		Host:       h.cfg.SMTP.Host,
		Port:       h.cfg.SMTP.Port,
		StartTLS:   h.cfg.SMTP.StartTLS,
		TimeoutSec: h.cfg.SMTP.TimeoutSec,
	}

	raw, err := smtppkg.Send(smtpCfg, s.Username, pass, msg)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}

	if conn, imapErr := imapConn(h.cfg, s); imapErr == nil {
		defer conn.Close()
		sentFolder := resolveFolderByIcon(conn, "sent", "Sent")
		_ = conn.AppendMessage(sentFolder, []goimap.Flag{goimap.FlagSeen}, raw)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
}

func (h *ComposeHandler) SaveDraft(c *echo.Context) error {
	s := c.Get("imap_session").(*session.IMAPSession)

	fromEmail := c.FormValue("from_email")
	if fromEmail == "" {
		fromEmail = s.Username
	}
	msg := &smtppkg.Message{
		From:        fromEmail,
		DisplayName: c.FormValue("from_name"),
		To:          splitAddrs(c.FormValue("to")),
		Cc:          splitAddrs(c.FormValue("cc")),
		Bcc:         splitAddrs(c.FormValue("bcc")),
		Subject:     c.FormValue("subject"),
		TextHTML:    c.FormValue("body_html"),
		TextPlain:   c.FormValue("body_plain"),
	}
	if replyTo := strings.TrimSpace(c.FormValue("reply_to")); replyTo != "" {
		msg.ReplyTo = replyTo
	}

	if err := h.attachFormFiles(c, msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	raw, err := smtppkg.Serialize(msg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer conn.Close()

	draftsFolder := resolveFolderByIcon(conn, "drafts", "Drafts")
	if err := conn.AppendMessage(draftsFolder, []goimap.Flag{goimap.FlagDraft}, raw); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok", "folder": draftsFolder})
}

func (h *ComposeHandler) UploadAttachment(c *echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file is required"})
	}

	maxBytes := int64(h.cfg.Upload.MaxSizeMB) * 1024 * 1024
	if maxBytes <= 0 {
		maxBytes = 25 * 1024 * 1024
	}
	if file.Size > maxBytes {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file too large"})
	}

	tempDir := h.cfg.Upload.TempDir
	if tempDir == "" {
		tempDir = "./tmp/uploads"
	}
	if err := os.MkdirAll(tempDir, 0o700); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "upload dir unavailable"})
	}

	id, err := randomID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "upload failed"})
	}

	dest := filepath.Join(tempDir, id)
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "upload failed"})
	}
	defer src.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "upload failed"})
	}
	defer out.Close()

	written, err := io.Copy(out, io.LimitReader(src, maxBytes+1))
	if err != nil || written > maxBytes {
		os.Remove(dest)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file too large"})
	}

	ct := file.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}

	return c.JSON(http.StatusOK, map[string]any{
		"files": []map[string]any{{
			"id":           id,
			"filename":     file.Filename,
			"size":         written,
			"content_type": ct,
		}},
	})
}

func (h *ComposeHandler) attachFormFiles(c *echo.Context, msg *smtppkg.Message) error {
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return nil
	}

	maxBytes := int64(h.cfg.Upload.MaxSizeMB) * 1024 * 1024
	if maxBytes <= 0 {
		maxBytes = 25 * 1024 * 1024
	}

	for _, file := range form.File["attachments"] {
		if file.Size > maxBytes {
			return fmt.Errorf("attachment too large")
		}
		src, err := file.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(io.LimitReader(src, maxBytes+1))
		src.Close()
		if err != nil || int64(len(data)) > maxBytes {
			return fmt.Errorf("attachment too large")
		}

		ct := file.Header.Get("Content-Type")
		if ct == "" {
			ct = "application/octet-stream"
		}
		msg.Attachments = append(msg.Attachments, smtppkg.Attachment{
			Filename:    file.Filename,
			ContentType: ct,
			Data:        data,
		})
	}
	return nil
}

func resolveFolderByIcon(conn *imap.Client, iconType, fallback string) string {
	if boxes, err := conn.ListMailboxes(); err == nil {
		for _, mb := range boxes {
			if mb.IconType == iconType {
				return mb.Name
			}
		}
	}
	return fallback
}

func splitAddrs(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
