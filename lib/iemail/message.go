package iemail

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"gopkg.in/gomail.v2"
)

type message struct {
	From    string
	To      string
	Subject string
	Date    string
	Body    string
}

func (m message) createRFC2822Message() (msgInRFC2822 string) {
	m.Date = time.Now().Local().Format(time.RFC1123)
	var templ *template.Template
	templ, _ = template.New("email").Parse(`
From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
Date: {{.Date}}
        
{{.Body}}
`)
	// 发生改动时注释以下代码，并接收template.New返回的error用于测试，测试完成后恢复
	var tpl bytes.Buffer
	templ.Execute(&tpl, m)
	msgInRFC2822 = strings.TrimSpace(tpl.String())

	//发生改动时以下代码解除注释
	// if err != nil {
	// 	err = fmt.Errorf("error parsing template: %w", err)
	// 	return
	// }

	// var tpl bytes.Buffer
	// if err = templ.Execute(&tpl, m); err != nil {
	// 	err = fmt.Errorf("error executing template: %w", err)
	// 	return
	// }
	// msgInRFC2822 = strings.TrimSpace(tpl.String())
	return
}

func (m message) createGoMailMessage() (msg *gomail.Message) {
	msg = gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.Body)
	msg.SetDateHeader("Date", time.Now())
	return
}
