package imapgoprebuilt

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

// get ticketmaster mfa code to join queue
func (n *ImapOpts) getWalmartMFA() (string, error) {
	// connect to server
	c, err := client.DialTLS(n.Imap.Imap, nil)
	if err != nil {
		return "", errors.New("could not connect to mail server")
	}

	// don't forget to logout
	defer c.Logout()

	// handle login
	if n.CatchallPass == "" {
		if err := c.Login(n.ReceiverEmail, n.ReceiverPass); err != nil {
			return "", fmt.Errorf("login and password are incorrect: %s:%s - %s", n.ReceiverEmail, n.ReceiverPass, err.Error())
		}
	} else {
		if err := c.Login(n.CatchallEmail, n.CatchallPass); err != nil {
			return "", fmt.Errorf("login and password are incorrect: %s:%s - %s", n.CatchallEmail, n.CatchallPass, err.Error())
		}
	}

	var boxes []string
	mailboxes := make(chan *imap.MailboxInfo, 5)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	for m := range mailboxes {
		if m.Name == "[Gmail]/Important" {
			continue
		}
		boxes = append(boxes, m.Name)
	}

	if err := <-done; err != nil {
		return "", fmt.Errorf("login and password are incorrect: %s:%s - %s", n.CatchallEmail, n.CatchallPass, err.Error())
	}

	var codesDate = make(map[string]int64)

	for _, box := range boxes {
		// Select INBOX
		mbox, err := c.Select(box, false)
		if err != nil {
			continue
		}

		// Get the last message
		if mbox.Messages == 0 {
			continue
		}

		var to, from uint32
		if mbox.Messages > 30 {
			from = mbox.Messages
			to = mbox.Messages - 30
		} else {
			from = mbox.Messages
			to = 0
		}

		seqSet := new(imap.SeqSet)
		seqSet.AddRange(from, to)

		// Get the whole message body
		var section imap.BodySectionName
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, 8)

		go func() {
			c.Fetch(seqSet, items, messages)
		}()

		var address, fromaddress, mailsubject string
		var maildate int64

		for msg := range messages {

			// If the message is null or if the activation email was found then skip the email
			if msg == nil {
				continue
			}

			r := msg.GetBody(&section)
			if r == nil {
				continue
			}

			// Create a new mail reader
			mr, err := mail.CreateReader(r)
			if err != nil {
				continue
			}

			// Print some info about the message
			header := mr.Header

			if date, err := header.Date(); err == nil {
				maildate = date.Unix()
			}

			if to, err := header.AddressList("To"); err == nil {
				if len(to) == 0 {
					continue
				}
				address = to[0].String()
			}

			if !strings.Contains(strings.ToLower(address), strings.ToLower(n.ReceiverEmail)) {
				continue
			}

			if from, err := header.AddressList("From"); err == nil {
				fromaddress = from[0].String()
			}

			if !strings.Contains(strings.ToLower(fromaddress), "walmart") {
				continue
			}
			if subject, err := header.Subject(); err == nil {
				mailsubject = strings.ToLower(subject)
			}

			if !strings.Contains(strings.ToLower(mailsubject), "code") {
				continue
			}

			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					break
				}

				switch p.Header.(type) {
				case *mail.InlineHeader:
					// This is the message's text (can be plain-text or HTML)
					b, _ := io.ReadAll(p.Body)

					re := regexp.MustCompile(`<strong>(\d{6})\s</strong>`)

					// Find the first match
					match := re.FindStringSubmatch(string(b))

					if len(match) > 1 {
						// The first capturing group contains the verification code
						verificationCode := match[1]
						codesDate[verificationCode] = maildate
					}
				default:
					continue
				}
			}
		}
	}

	var code string
	var time int64
	for k, v := range codesDate {
		if v > time {
			time = v
			code = k
		}
	}
	if code == "" {
		return code, errors.New("email not yet found")
	} else {
		return code, nil
	}
}
