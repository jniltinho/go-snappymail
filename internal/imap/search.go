package imap

import "github.com/emersion/go-imap/v2"

// SearchCriteria defines filter parameters for an IMAP search request.
type SearchCriteria struct {
	Subject string
	From    string
	To      string
	Body    string
	Unseen  bool
	Flagged bool
}

// Search executes a UID search on the currently selected mailbox and returns matching UIDs.
func (c *Client) Search(criteria *SearchCriteria) ([]imap.UID, error) {
	sc := &imap.SearchCriteria{}

	var textCriteria []imap.SearchCriteria
	if criteria.Subject != "" {
		textCriteria = append(textCriteria, imap.SearchCriteria{
			Header: []imap.SearchCriteriaHeaderField{{Key: "Subject", Value: criteria.Subject}},
		})
	}
	if criteria.From != "" {
		textCriteria = append(textCriteria, imap.SearchCriteria{
			Header: []imap.SearchCriteriaHeaderField{{Key: "From", Value: criteria.From}},
		})
	}
	if criteria.To != "" {
		textCriteria = append(textCriteria, imap.SearchCriteria{
			Header: []imap.SearchCriteriaHeaderField{{Key: "To", Value: criteria.To}},
		})
	}
	if criteria.Body != "" {
		textCriteria = append(textCriteria, imap.SearchCriteria{Body: []string{criteria.Body}})
	}

	for len(textCriteria) > 1 {
		last := textCriteria[len(textCriteria)-1]
		prev := textCriteria[len(textCriteria)-2]
		merged := imap.SearchCriteria{}
		merged.Or = append(merged.Or, [2]imap.SearchCriteria{prev, last})
		textCriteria = append(textCriteria[:len(textCriteria)-2], merged)
	}
	if len(textCriteria) == 1 {
		sc = &textCriteria[0]
	}

	if criteria.Unseen {
		sc.NotFlag = append(sc.NotFlag, imap.FlagSeen)
	}
	if criteria.Flagged {
		sc.Flag = append(sc.Flag, imap.FlagFlagged)
	}

	data, err := c.Client.UIDSearch(sc, nil).Wait()
	if err != nil {
		return nil, err
	}
	return data.AllUIDs(), nil
}
