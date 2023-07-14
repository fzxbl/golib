package iemail

import (
	"fmt"
	"testing"
)

func Test_Send(t *testing.T) {
	var cli Client
	cli.MustInit("./testdata/email.toml", ClientModeReceiveAndSend)
	msg := Message{
		From:    cli.cfg.Auth.User,
		To:      "xxx@xxmail.com",
		Subject: "This is a test",
		Body:    "hello, world",
	}
	err := cli.Send(msg, true)
	fmt.Println(err)
}

func Test_MustInitClient(t *testing.T) {
	var cli Client
	cli.MustInit("./testdata/email.toml", ClientModeReceiveAndSend)
	fmt.Println(cli)
}
