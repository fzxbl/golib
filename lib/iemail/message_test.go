package iemail

import (
	"fmt"
	"testing"
)

func Test_CreateEmail(t *testing.T) {
	m := Message{
		From:    "alice@example.com",
		To:      "bob@example.com",
		Subject: "Hello",
		Body:    "Hello, I hope you are doing well.",
	}
	msg := m.CreateRFC2822Message()

	fmt.Println(msg)
}
