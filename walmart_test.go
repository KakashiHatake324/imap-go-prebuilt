package imapgoprebuilt

import (
	"log"
	"testing"
)

func TestWalmart(t *testing.T) {
	imapOpts := &ImapOpts{
		Imap:          Gmail,
		Site:          Walmart,
		ReceiverEmail: "rafaeltorres324@gmail.com",
		ReceiverPass:  "umrjlnhtaathgyhg",
		CatchallEmail: "",
		CatchallPass:  "",
		MaxChecks:     5,
	}
	code, err := imapOpts.FetchEmail()
	log.Println(code, err)
}
