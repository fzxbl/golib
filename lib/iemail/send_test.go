package iemail

import (
	"fmt"
	"testing"
)

func Test_Send(t *testing.T) {
	cli := NewClient("./testdata/email.toml", WithMode(ClientModeSend))
	msg := message{
		From:    cli.Config.Auth.User,
		To:      "xxx@xxmail.com",
		Subject: "This is a test",
		Body:    "hello, world",
	}
	err := cli.send(msg, true)
	fmt.Println(err)
}
