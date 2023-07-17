package irequest

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_Client(t *testing.T) {
	i := InitCookie{Host: "http://www.baidu.com", Domain: ".baidu.com", CookieHeader: `a=b=2;c=3`}
	jari := i.initCookie()
	urlobj, _ := url.Parse("http://api.scopus.com.zzulib.vpn358.com/api/documents/search")
	cookiesi := jari.Cookies(urlobj)

	fmt.Println(cookiesi)
}
