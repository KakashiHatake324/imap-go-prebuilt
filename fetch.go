package imapgoprebuilt

import (
	"time"
)

// fetch the email information with the prebuilt functions
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
	case TicketMaster:
		for i := 1; i < n.MaxChecks; i++ {
			message, err = n.getTicketMasterMFA()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	case Walmart:
		for i := 1; i < n.MaxChecks; i++ {
			message, err = n.getWalmartMFA()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	case Footsites:
		for i := 1; i < n.MaxChecks; i++ {
			message, err = n.getFLXActivationLink()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
	return message, err
}
