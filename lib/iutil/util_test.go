package iutil

import (
	"fmt"
	"testing"
)

func Test_BlockIfExpired(t *testing.T) {
	BlockIfExpired(2023, 06, 26, 14)
}

func Test_TemplateReplace(t *testing.T) {
	temp := "Hello, my name is {{.Name}} and I am from {{.Country}}."
	data := map[string]string{
		"Name":    "Bob",
		"Country": "United States",
	}

	r1, e1 := TemplateReplace(temp, data)
	fmt.Println(r1, e1)
	type tmp struct {
		Name    string
		Country string
	}
	data2 := tmp{
		Name:    "ALICE",
		Country: "China",
	}
	r2, e2 := TemplateReplace(temp, data2)
	fmt.Println(r2, e2)
}
