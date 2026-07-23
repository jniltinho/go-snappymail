package handler

import (
	"fmt"
	"strings"

	"go-snappymail/internal/config"
	imappkg "go-snappymail/internal/imap"
)

// MessageHandler handles reading, flagging, moving, and deleting individual email messages.
type MessageHandler struct {
	cfg *config.Config
}

func messageDownloadName(subject string, uid uint64) string {
	name := strings.TrimSpace(subject)
	if name == "" {
		return fmt.Sprintf("message-%d.eml", uid)
	}

	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "-",
	)
	name = replacer.Replace(name)
	name = strings.Join(strings.Fields(name), " ")
	if name == "" {
		return fmt.Sprintf("message-%d.eml", uid)
	}
	return name + ".eml"
}

func findTrashFolder(conn *imappkg.Client) string {
	trashFolder := "Trash"
	if boxes, err := conn.ListMailboxes(); err == nil {
		for _, mb := range boxes {
			if mb.IsTrash {
				return mb.Name
			}
		}
	}
	return trashFolder
}
