package iutil

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/fzxbl/golib/lib/interact"
)

// BlockIfExpired 如果过期,则阻塞,等待用户输入Ctrl+C,参数为过期时间
func BlockIfExpired(year, month, day, hour int) {
	shanghai := time.FixedZone("Asia/Shanghai", 8*60*60)
	expiredTime := time.Date(year, time.Month(month), day, hour, 0, 0, 0, shanghai)
	log.Printf("程序在%s前可用", expiredTime.Format("2006-01-02 15:04:05"))
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		log.Print("获取时间错误，请重新执行程序", err)
		interact.BlockOnSignal()
	}
	defer resp.Body.Close()
	bdTime := resp.Header.Get("Date")
	now, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", bdTime)
	if err != nil {
		log.Print("请联系开发者", err)
		interact.BlockOnSignal()
	}
	if now.After(expiredTime) {
		log.Print("时间过期，请联系开发者")
		interact.BlockOnSignal()
	}
}

func TemplateReplace(temp string, data any) (result string, err error) {
	tmpl, err := template.New("test").Parse(temp)
	if err != nil {
		return
	}
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, data)
	if err != nil {
		return
	}
	result = tmplBytes.String()
	return
}
