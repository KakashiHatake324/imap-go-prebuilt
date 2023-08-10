package imapgoprebuilt

import (
	"time"
)

func (n *ImapOpts) FetchEmail() (string, error) {
	var message string
	var err error

	switch n.Site {
	case Nike:
		for i := 1; i < n.MaxChecks; i++ {
			message, err = n.getNikeLoginCode()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
	return message, err
}
