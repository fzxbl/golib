package iemail

import (
	"crypto/tls"
	"strconv"

	"github.com/emersion/go-imap/v2/imapclient"
	iconf "github.com/fzxbl/golib/conf"
	"gopkg.in/gomail.v2"
)

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

func MustInitClient(filepath string) (cli Client) {
	var cfg Config
	iconf.MustParseToml(filepath, &cfg)
	smtpCli, imapCli := mustInitSmtpImap(cfg)
	cli.cfg = cfg
	cli.smtpClient = smtpCli
	cli.imapClient = imapCli
	return
}

// mustInitSmtp 初始化SMTP服务，使用gomail.NewMessage()构造msg,使用client.DialAndSend(msg)发送
func mustInitSmtp(cfg Config) (client *gomail.Dialer) {
	client = gomail.NewDialer(cfg.SendServer.Adress, cfg.SendServer.Port, cfg.Auth.User, cfg.Auth.Passwd)
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return
}

// mustInitImap 初始化Imap Client, 使用完成后务必调用client.Logout()
func mustInitImap(cfg Config) (client *imapclient.Client) {
	var err error
	client, err = imapclient.DialTLS(cfg.ReceiveServer.Adress+":"+strconv.Itoa(cfg.ReceiveServer.Port), nil)
	if err != nil {
		panic(err)
	}
	if err := client.Login(cfg.Auth.User, cfg.Auth.Passwd).Wait(); err != nil {
		panic(err)
	}
	return
}

// mustInitSmtpImap 同时初始化SMTP和Imap服务，使用完成后务必调用imapClient.Logout()
func mustInitSmtpImap(cfg Config) (smtpClient *gomail.Dialer, imapClient *imapclient.Client) {
	smtpClient = gomail.NewDialer(cfg.SendServer.Adress, cfg.SendServer.Port, cfg.Auth.User, cfg.Auth.Passwd)
	smtpClient.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	var err error
	imapClient, err = imapclient.DialTLS(cfg.ReceiveServer.Adress+":"+strconv.Itoa(cfg.ReceiveServer.Port), nil)
	if err != nil {
		panic(err)
	}
	if err := imapClient.Login(cfg.Auth.User, cfg.Auth.Passwd).Wait(); err != nil {
		panic(err)
	}
	return
}
