package iemail

import (
	"fmt"
	"testing"
)

func Test_Send(t *testing.T) {
	cli := MustInitClient("./testdata/email.toml")
	msg := Message{
		From:    cli.cfg.Auth.User,
		To:      "luckygj@foxmail.com",
		Subject: "This is a test",
		Body:    "hello, world",
	}
	err := cli.Send(msg)
	fmt.Println(err)
}

func Test_MustInitClient(t *testing.T) {
	cli := MustInitClient("./testdata/email.toml")
	fmt.Println(cli)
}
