package imap

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"strings"
	"time"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
)

// ParsedMessage holds the decoded body parts and attachments of a MIME message.
type ParsedMessage struct {
	TextPlain    string
	TextHTML     string
	Attachments  []Attachment
	CalendarInfo *CalendarInfo
}

// CalendarInfo holds iCalendar event data (RFC 5545).
type CalendarInfo struct {
	Method    string   `json:"method"`
	Summary   string   `json:"summary"`
	Organizer string   `json:"organizer"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Location  string   `json:"location"`
	UID       string   `json:"uid"`
	Attendees []string `json:"attendees"`
}

// Attachment describes a single file attached to an email message.
type Attachment struct {
	Filename    string
	ContentType string
	Size        int64
	Part        int
	Data        []byte
}

func formatICalDate(s string) string {
	s = strings.TrimRight(s, "Z")
	for _, layout := range []string{"20060102T150405", "20060102T1504", "20060102"} {
		if t, err := time.Parse(layout, s); err == nil {
			if layout == "20060102" {
				return t.Format("02/01/2006")
			}
			return t.Format("02/01/2006 15:04")
		}
	}
	return s
}

func parseICS(data []byte) *CalendarInfo {
	info := &CalendarInfo{}
	inVEvent := false
	for raw := range strings.SplitSeq(string(data), "\n") {
		line := strings.TrimRight(raw, "\r")
		upper := strings.ToUpper(strings.TrimSpace(line))
		if upper == "BEGIN:VEVENT" {
			inVEvent = true
			continue
		}
		if upper == "END:VEVENT" {
			inVEvent = false
			continue
		}

		keyRaw, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key := strings.ToUpper(strings.SplitN(keyRaw, ";", 2)[0])
		val = strings.TrimSpace(val)

		if !inVEvent {
			if key == "METHOD" {
				info.Method = strings.ToUpper(val)
			}
			continue
		}
		switch key {
		case "SUMMARY":
			info.Summary = val
		case "ORGANIZER":
			info.Organizer = strings.TrimPrefix(strings.ToLower(val), "mailto:")
		case "DTSTART":
			info.StartTime = formatICalDate(val)
		case "DTEND":
			info.EndTime = formatICalDate(val)
		case "LOCATION":
			info.Location = val
		case "UID":
			info.UID = val
		case "ATTENDEE":
			email := strings.TrimPrefix(strings.ToLower(val), "mailto:")
			if email != "" {
				info.Attendees = append(info.Attendees, email)
			}
		}
	}
	return info
}

// ParseMessage decodes a raw RFC822 message into its text bodies and attachments.
func ParseMessage(raw []byte) (*ParsedMessage, error) {
	e, err := message.Read(bytes.NewReader(raw))
	if err != nil && !message.IsUnknownCharset(err) && !message.IsUnknownEncoding(err) {
		return rawFallback(raw), nil
	}
	if e == nil {
		return rawFallback(raw), nil
	}

	pm := &ParsedMessage{}
	cidMap := make(map[string]string)
	partNum := 0

	e.Walk(func(_ []int, entity *message.Entity, walkErr error) error {
		if walkErr != nil && !message.IsUnknownCharset(walkErr) && !message.IsUnknownEncoding(walkErr) {
			return nil
		}
		if entity == nil {
			return nil
		}
		if entity.MultipartReader() != nil {
			return nil
		}

		partNum++
		processLeaf(entity, pm, partNum, cidMap)
		return nil
	})

	if pm.TextHTML != "" && len(cidMap) > 0 {
		for cid, dataURI := range cidMap {
			pm.TextHTML = strings.ReplaceAll(pm.TextHTML, "cid:"+cid, dataURI)
		}
	}

	if pm.TextPlain == "" && pm.TextHTML == "" && pm.CalendarInfo == nil {
		return rawFallback(raw), nil
	}

	return pm, nil
}

func processLeaf(e *message.Entity, pm *ParsedMessage, partNum int, cidMap map[string]string) {
	ctRaw := e.Header.Get("Content-Type")
	dispRaw := e.Header.Get("Content-Disposition")
	contentID := strings.Trim(e.Header.Get("Content-ID"), "<>")

	mediaType, ctParams, _ := mime.ParseMediaType(ctRaw)
	dispType, dispParams, _ := mime.ParseMediaType(dispRaw)

	data, err := io.ReadAll(e.Body)
	if err != nil {
		return
	}

	filename := dispParams["filename"]
	if filename == "" {
		filename = ctParams["name"]
	}

	switch {
	case dispType == "attachment" || (filename != "" && dispType != "inline"):
		pm.Attachments = append(pm.Attachments, Attachment{
			Filename:    filename,
			ContentType: mediaType,
			Size:        int64(len(data)),
			Part:        partNum,
			Data:        data,
		})

	case contentID != "" && strings.HasPrefix(mediaType, "image/"):
		b64 := base64.StdEncoding.EncodeToString(data)
		cidMap[contentID] = fmt.Sprintf("data:%s;base64,%s", mediaType, b64)

	case strings.HasPrefix(mediaType, "text/html"):
		if pm.TextHTML == "" {
			pm.TextHTML = string(data)
		}

	case strings.HasPrefix(mediaType, "text/plain"):
		if pm.TextPlain == "" {
			pm.TextPlain = strings.TrimSpace(string(data))
		}

	case strings.HasPrefix(mediaType, "text/calendar"),
		strings.HasSuffix(mediaType, "ics"),
		strings.EqualFold(mediaType, "application/ics"):
		pm.CalendarInfo = parseICS(data)
	}
}

func rawFallback(raw []byte) *ParsedMessage {
	pm := &ParsedMessage{}
	if _, body, ok := bytes.Cut(raw, []byte("\r\n\r\n")); ok {
		pm.TextPlain = strings.TrimSpace(string(body))
	} else if _, body, ok := bytes.Cut(raw, []byte("\n\n")); ok {
		pm.TextPlain = strings.TrimSpace(string(body))
	}
	return pm
}
