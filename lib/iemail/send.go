package iemail

import "fmt"

func (c Client) Send(to string, subject string, content string, attachFileName string, syncSentBox bool) (err error) {
	m := message{
		From:    c.Config.Auth.User,
		To:      to,
		Subject: subject,
		Body:    content,
	}
	err = c.send(m, attachFileName, syncSentBox)
	return
}

// send syncSentBox是否同步到发件箱
func (c Client) send(m message, attachFileName string, syncSentBox bool) (err error) {
	sendMsg := m.createGoMailMessage()
	if attachFileName != "" {
		sendMsg.Attach(attachFileName)
	}
	err = c.smtpClient.DialAndSend(sendMsg)
	if err != nil || !syncSentBox || c.imapClient == nil {
		return
	}
	defer c.imapClient.Logout()
	archivedMsg := m.createRFC2822Message()
	data := []byte(archivedMsg)
	appendCmd := c.imapClient.Append(c.Config.MailBox.SentBox, int64(len(data)), nil)
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
