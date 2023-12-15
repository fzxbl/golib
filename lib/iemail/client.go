package iemail

import (
	"crypto/tls"
	"strconv"

	"github.com/emersion/go-imap/v2/imapclient"
	iconf "github.com/fzxbl/golib/iconf"
	"gopkg.in/gomail.v2"
)

type Client struct {
	smtpClient *gomail.Dialer
	imapClient *imapclient.Client
	Config     Config
}

type Options struct {
	ClientMode ClientMode
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

type Option func(opts *Options)

func WithMode(mode ClientMode) Option {
	return func(opts *Options) {
		opts.ClientMode = mode
	}
}

func NewClient(confPath string, options ...Option) *Client {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	cli := &Client{}
	var cfg Config
	iconf.MustParseToml(confPath, &cfg)
	cli.Config = cfg

	if opts.ClientMode == 0 {
		imapCli := mustInitImap(cfg)
		cli.imapClient = imapCli
		smtpClient := mustInitSmtp(cfg)
		cli.smtpClient = smtpClient
	}
	if opts.ClientMode&1 == 1 {
		imapCli := mustInitImap(cfg)
		cli.imapClient = imapCli
	}
	if opts.ClientMode&2 == 2 {
		smtpClient := mustInitSmtp(cfg)
		cli.smtpClient = smtpClient
	}
	return cli
}

func (cli *Client) Close() {
	if cli.imapClient != nil {
		cli.imapClient.Logout()
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
