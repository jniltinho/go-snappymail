package handler

import (
	"net/http"
	"strconv"

	"go-snappymail/internal/session"

	"github.com/emersion/go-imap/v2"
	"github.com/labstack/echo/v5"
)

func (h *MessageHandler) Flag(c *echo.Context) error {
	mailbox := c.Param("mailbox")
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 32)
	if err != nil {
		return echo.ErrBadRequest
	}

	flag := c.FormValue("flag")
	value := c.FormValue("value") == "1"

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SelectMailbox(mailbox); err != nil {
		return echo.ErrNotFound
	}

	imapUID := imap.UID(uid)
	switch flag {
	case "seen":
		if value {
			err = conn.MarkSeen(imapUID)
		} else {
			err = conn.MarkUnseen(imapUID)
		}
	case "flagged":
		err = conn.MarkFlagged(imapUID, value)
	case "answered":
		err = conn.MarkAnswered(imapUID, value)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid flag"})
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MessageHandler) Move(c *echo.Context) error {
	mailbox := c.Param("mailbox")
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 32)
	if err != nil {
		return echo.ErrBadRequest
	}
	dest := c.FormValue("dest")
	if dest == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dest is required"})
	}

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SelectMailbox(mailbox); err != nil {
		return echo.ErrNotFound
	}
	if err := conn.MoveMessage(imap.UID(uid), dest); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MessageHandler) Delete(c *echo.Context) error {
	mailbox := c.Param("mailbox")
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 32)
	if err != nil {
		return echo.ErrBadRequest
	}

	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SelectMailbox(mailbox); err != nil {
		return echo.ErrNotFound
	}

	trashFolder := findTrashFolder(conn)

	if mailbox == trashFolder {
		err = conn.DeleteMessage(imap.UID(uid))
	} else {
		err = conn.MoveMessage(imap.UID(uid), trashFolder)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MessageHandler) EmptyTrash(c *echo.Context) error {
	mailbox := c.Param("mailbox")
	s := c.Get("imap_session").(*session.IMAPSession)
	conn, err := imapConn(h.cfg, s)
	if err != nil {
		return err
	}
	defer conn.Close()

	trashFolder := findTrashFolder(conn)

	if mailbox == trashFolder {
		if err := conn.EmptyMailbox(mailbox); err != nil {
			return err
		}
	} else {
		if err := conn.MoveAllMessages(mailbox, trashFolder); err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
