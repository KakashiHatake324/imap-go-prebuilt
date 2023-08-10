package imapgoprebuilt

type ImapOpts struct {
	Imap          *EmailOpts
	Site          string
	ReceiverEmail string
	ReceiverPass  string
	CatchallEmail string
	CatchallPass  string
	MaxChecks     int
}

type EmailOpts struct {
	Email string
	Imap  string
}
