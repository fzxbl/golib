package iemail

import (
	"crypto/tls"
	"strconv"

	"github.com/emersion/go-imap/v2/imapclient"
	iconf "github.com/fzxbl/golib/iconf"
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
type ClientMode uint8

const (
	// 仅接收邮件
	ClientModeReceive ClientMode = iota + 1
	// 仅发送邮件，无法同步发件箱
	ClientModeSend ClientMode = iota + 1
	// 发送和接收邮件，可同步发件箱
	ClientModeReceiveAndSend ClientMode = iota + 1
)

func (cli *Client) MustInit(confPath string, mode ClientMode) {
	var cfg Config
	iconf.MustParseToml(confPath, &cfg)
	cli.cfg = cfg
	if mode&1 == 1 {
		imapCli := mustInitImap(cfg)
		cli.imapClient = imapCli
	}
	if mode&2 == 1 {
		smtpClient := mustInitSmtp(cfg)
		cli.smtpClient = smtpClient
	}
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
