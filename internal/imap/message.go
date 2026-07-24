package imap

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap/v2"
)

// Envelope contains the message headers needed to render a mailbox listing row.
type Envelope struct {
	UID       imap.UID `json:"uid"`
	Subject   string   `json:"subject"`
	From      string   `json:"from"`
	FromEmail string   `json:"from_email"`
	To        string   `json:"to"`
	Date      string   `json:"date"`
	DateTS    int64    `json:"date_ts"`
	HasAttach bool     `json:"has_attachment"`
	Seen      bool     `json:"seen"`
	Flagged   bool     `json:"flagged"`
	Answered  bool     `json:"answered"`
	Size      int64    `json:"size"`
}

// FetchEnvelopes retrieves the envelope (headers and flags) for a slice of UIDs.
// The mailbox must be selected before calling this method.
func (c *Client) FetchEnvelopes(uids []imap.UID) ([]Envelope, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	seqSet := imap.UIDSetNum(uids...)
	msgs, err := c.Client.Fetch(seqSet, &imap.FetchOptions{
		UID:           true,
		Flags:         true,
		Envelope:      true,
		RFC822Size:    true,
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
	}).Collect()
	if err != nil {
		return nil, err
	}

	result := make([]Envelope, 0, len(msgs))
	for _, m := range msgs {
		env := Envelope{UID: m.UID}
		if m.Envelope != nil {
			env.Subject = m.Envelope.Subject
			if len(m.Envelope.From) > 0 {
				addr := m.Envelope.From[0]
				env.FromEmail = addr.Mailbox + "@" + addr.Host
				if addr.Name != "" && addr.Name != env.FromEmail {
					env.From = addr.Name
				} else {
					env.From = env.FromEmail
				}
			}
			toAddrs := make([]string, 0, len(m.Envelope.To))
			for _, addr := range m.Envelope.To {
				email := addr.Mailbox + "@" + addr.Host
				if addr.Name != "" && addr.Name != email {
					toAddrs = append(toAddrs, addr.Name+" <"+email+">")
				} else {
					toAddrs = append(toAddrs, email)
				}
			}
			env.To = strings.Join(toAddrs, ", ")
			env.Date = m.Envelope.Date.Format("02/01/2006 15:04")
			if !m.Envelope.Date.IsZero() {
				env.DateTS = m.Envelope.Date.Unix()
			}
		}
		env.Size = m.RFC822Size
		if m.BodyStructure != nil {
			m.BodyStructure.Walk(func(_ []int, part imap.BodyStructure) bool {
				if d := part.Disposition(); d != nil && strings.EqualFold(d.Value, "attachment") {
					env.HasAttach = true
					return false
				}
				return true
			})
		}
		for _, f := range m.Flags {
			switch f {
			case imap.FlagSeen:
				env.Seen = true
			case imap.FlagFlagged:
				env.Flagged = true
			case imap.FlagAnswered:
				env.Answered = true
			}
		}
		result = append(result, env)
	}
	return result, nil
}

// FetchRawMessage downloads the complete RFC822 bytes of a message using BODY.PEEK[].
func (c *Client) FetchRawMessage(uid imap.UID) ([]byte, error) {
	seqSet := imap.UIDSetNum(uid)
	section := &imap.FetchItemBodySection{Peek: true}

	msgs, err := c.Client.Fetch(seqSet, &imap.FetchOptions{
		UID:         true,
		BodySection: []*imap.FetchItemBodySection{section},
	}).Collect()
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	return msgs[0].FindBodySection(section), nil
}

// MarkSeen sets the \Seen flag on the message with the given UID.
func (c *Client) MarkSeen(uid imap.UID) error {
	return c.Client.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
		Op:    imap.StoreFlagsAdd,
		Flags: []imap.Flag{imap.FlagSeen},
	}, nil).Close()
}

// MarkUnseen removes the \Seen flag from the message with the given UID.
func (c *Client) MarkUnseen(uid imap.UID) error {
	return c.Client.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
		Op:    imap.StoreFlagsDel,
		Flags: []imap.Flag{imap.FlagSeen},
	}, nil).Close()
}

// MarkFlagged adds or removes the \Flagged flag on the message with the given UID.
func (c *Client) MarkFlagged(uid imap.UID, flagged bool) error {
	op := imap.StoreFlagsAdd
	if !flagged {
		op = imap.StoreFlagsDel
	}
	return c.Client.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
		Op:    op,
		Flags: []imap.Flag{imap.FlagFlagged},
	}, nil).Close()
}

// MarkAnswered adds or removes the \Answered flag on the message with the given UID.
func (c *Client) MarkAnswered(uid imap.UID, answered bool) error {
	op := imap.StoreFlagsAdd
	if !answered {
		op = imap.StoreFlagsDel
	}
	return c.Client.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
		Op:    op,
		Flags: []imap.Flag{imap.FlagAnswered},
	}, nil).Close()
}

// MoveMessage moves a single message to the named destination folder using IMAP MOVE.
func (c *Client) MoveMessage(uid imap.UID, dest string) error {
	_, err := c.Client.Move(imap.UIDSetNum(uid), dest).Wait()
	return err
}

// DeleteMessage permanently removes a message by setting \Deleted and running EXPUNGE.
func (c *Client) DeleteMessage(uid imap.UID) error {
	if err := c.Client.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
		Op:    imap.StoreFlagsAdd,
		Flags: []imap.Flag{imap.FlagDeleted},
	}, nil).Close(); err != nil {
		return err
	}
	return c.Client.Expunge().Close()
}
