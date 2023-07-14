package iemail

import (
	"fmt"

	"github.com/emersion/go-imap/v2/imapclient"
	"gopkg.in/gomail.v2"
)

type ClientIter interface {
	Send() error
}

type Client struct {
	smtpClient *gomail.Dialer
	imapClient *imapclient.Client
	cfg        Config
}

func (c Client) Send(m Message) (err error) {
	defer c.imapClient.Logout()
	msg1 := m.CreateGoMailMessage()
	err = c.smtpClient.DialAndSend(msg1)
	if err != nil {
		return
	}
	msg2 := m.CreateRFC2822Message()
	data := []byte(msg2)

	appendCmd := c.imapClient.Append(c.cfg.MailBox.SentBox, int64(len(data)), nil)
	if _, err = appendCmd.Write(data); err != nil {
		err = fmt.Errorf("failed to write message: %w", err)
		return
	}
	if err = appendCmd.Close(); err != nil {
		err = fmt.Errorf("failed to close append cmd: %w", err)
		return
	}
	if _, err = appendCmd.Wait(); err != nil {
		err = fmt.Errorf("failed to wait append cmd: %w", err)
		return
	}
	return
}
