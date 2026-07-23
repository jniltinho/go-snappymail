// Package smtp provides email composition and delivery via SMTP.
package smtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	netmail "net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/wneessen/go-mail"
)

// Message holds all fields needed to compose an outgoing email.
type Message struct {
	From        string
	DisplayName string
	ReplyTo     string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	TextPlain   string
	TextHTML    string
	Attachments []Attachment
}

// Attachment holds the binary content of a file to attach to an outgoing message.
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// Config holds SMTP server connection settings.
type Config struct {
	Host       string
	Port       int
	StartTLS   bool
	TimeoutSec int
}

var dataURIRe = regexp.MustCompile(`src="data:([^;]+);base64,([^"]+)"`)

type inlineImage struct {
	cid      string
	filename string
	mimeType string
	data     []byte
}

func extractInlineImages(html string) (string, []inlineImage) {
	var inlines []inlineImage
	idx := 0
	result := dataURIRe.ReplaceAllStringFunc(html, func(match string) string {
		parts := dataURIRe.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		mimeType := parts[1]
		b64data := strings.ReplaceAll(strings.ReplaceAll(parts[2], "\n", ""), " ", "")

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			return match
		}

		idx++
		cid := fmt.Sprintf("inline-img-%d", idx)
		ext := "bin"
		switch mimeType {
		case "image/png":
			ext = "png"
		case "image/jpeg", "image/jpg":
			ext = "jpg"
		case "image/gif":
			ext = "gif"
		case "image/webp":
			ext = "webp"
		}
		filename := fmt.Sprintf("inline-%d.%s", idx, ext)

		inlines = append(inlines, inlineImage{cid: cid, filename: filename, mimeType: mimeType, data: data})
		return fmt.Sprintf(`src="cid:%s"`, cid)
	})
	return result, inlines
}

func buildMessage(msg *Message) (*mail.Msg, error) {
	m := mail.NewMsg()

	if msg.DisplayName != "" {
		if err := m.FromFormat(msg.DisplayName, msg.From); err != nil {
			return nil, fmt.Errorf("set from address: %w", err)
		}
	} else {
		if err := m.From(msg.From); err != nil {
			return nil, fmt.Errorf("set from address: %w", err)
		}
	}

	for _, toStr := range msg.To {
		addr, err := netmail.ParseAddress(toStr)
		if err != nil {
			if err := m.AddTo(toStr); err != nil {
				return nil, fmt.Errorf("invalid To address %s: %w", toStr, err)
			}
		} else {
			if err := m.AddToFormat(addr.Name, addr.Address); err != nil {
				return nil, fmt.Errorf("invalid To address %s: %w", toStr, err)
			}
		}
	}

	for _, ccStr := range msg.Cc {
		addr, err := netmail.ParseAddress(ccStr)
		if err != nil {
			if err := m.AddCc(ccStr); err != nil {
				return nil, fmt.Errorf("invalid Cc address %s: %w", ccStr, err)
			}
		} else {
			if err := m.AddCcFormat(addr.Name, addr.Address); err != nil {
				return nil, fmt.Errorf("invalid Cc address %s: %w", ccStr, err)
			}
		}
	}

	for _, bccStr := range msg.Bcc {
		addr, err := netmail.ParseAddress(bccStr)
		if err != nil {
			if err := m.AddBcc(bccStr); err != nil {
				return nil, fmt.Errorf("invalid Bcc address %s: %w", bccStr, err)
			}
		} else {
			if err := m.AddBccFormat(addr.Name, addr.Address); err != nil {
				return nil, fmt.Errorf("invalid Bcc address %s: %w", bccStr, err)
			}
		}
	}

	if msg.ReplyTo != "" {
		_ = m.ReplyTo(msg.ReplyTo)
	}

	m.Subject(msg.Subject)

	htmlBody := msg.TextHTML
	var inlines []inlineImage
	if htmlBody != "" {
		htmlBody, inlines = extractInlineImages(htmlBody)
	}

	if msg.TextPlain != "" && htmlBody != "" {
		m.SetBodyString(mail.TypeTextPlain, msg.TextPlain)
		m.AddAlternativeString(mail.TypeTextHTML, htmlBody)
	} else if msg.TextPlain != "" {
		m.SetBodyString(mail.TypeTextPlain, msg.TextPlain)
	} else if htmlBody != "" {
		m.SetBodyString(mail.TypeTextHTML, htmlBody)
	}

	for _, img := range inlines {
		m.EmbedReader(img.filename, bytes.NewReader(img.data),
			mail.WithFileContentID(img.cid),
			mail.WithFileContentType(mail.ContentType(img.mimeType)))
	}

	for _, a := range msg.Attachments {
		m.AttachReader(a.Filename, bytes.NewReader(a.Data),
			mail.WithFileContentType(mail.ContentType(a.ContentType)))
	}

	return m, nil
}

// Serialize builds RFC822 bytes for a message without sending it.
func Serialize(msg *Message) ([]byte, error) {
	m, err := buildMessage(msg)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err := m.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("serialize email: %w", err)
	}
	return buf.Bytes(), nil
}

// Send delivers msg via SMTP and returns the raw RFC822 bytes.
func Send(cfg Config, user, pass string, msg *Message) ([]byte, error) {
	m, err := buildMessage(msg)
	if err != nil {
		return nil, err
	}

	var opts []mail.Option
	opts = append(opts, mail.WithPort(cfg.Port))
	opts = append(opts, mail.WithTimeout(time.Duration(cfg.TimeoutSec)*time.Second))

	if user != "" || pass != "" {
		opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
		opts = append(opts, mail.WithUsername(user))
		opts = append(opts, mail.WithPassword(pass))
	}

	if cfg.Port == 465 {
		opts = append(opts, mail.WithSSL())
	} else if cfg.StartTLS {
		opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	} else {
		opts = append(opts, mail.WithTLSPolicy(mail.TLSOpportunistic))
	}

	client, err := mail.NewClient(cfg.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("initialize smtp client: %w", err)
	}

	var buf bytes.Buffer
	if _, err := m.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("serialize email: %w", err)
	}
	raw := buf.Bytes()

	if err := client.DialAndSend(m); err != nil {
		return nil, fmt.Errorf("send email: %w", err)
	}

	return raw, nil
}
