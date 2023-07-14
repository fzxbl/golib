package iemail

import (
	"fmt"

	"github.com/emersion/go-imap/v2/imapclient"
	"gopkg.in/gomail.v2"
)

type Client struct {
	smtpClient *gomail.Dialer
	imapClient *imapclient.Client
	cfg        Config
}

// Send syncSentBox是否同步到发件箱
func (c Client) Send(m Message, syncSentBox bool) (err error) {
	sendMsg := m.CreateGoMailMessage()
	err = c.smtpClient.DialAndSend(sendMsg)
	if err != nil || !syncSentBox || c.imapClient == nil {
		return
	}
	defer c.imapClient.Logout()
	archivedMsg := m.CreateRFC2822Message()
	data := []byte(archivedMsg)
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
