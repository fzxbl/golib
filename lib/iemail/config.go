package iemail

type Server struct {
	Adress string
	Port   int
	UseTLS bool
}

type Auth struct {
	User   string
	Passwd string
}

type MailBox struct {
	SentBox string
	InBox   string
}

type Config struct {
	SendServer    Server
	ReceiveServer Server
	Auth          Auth
	MailBox       MailBox
}
