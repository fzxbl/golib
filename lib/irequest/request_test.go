package irequest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_Request(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"message": "Hello, World!",
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
	fmt.Println("开始监听端口: 8080")
	go func() {
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
	}()
	time.Sleep(time.Second * 2)
	type exampleResp struct {
		Message string `json:"message"`
	}
	var data exampleResp
	c := NewClient(WithCookie(&InitCookie{Host: "http://localhost:8080", Domain: ".localhost:8080", CookieHeader: "a=b=1;b=2"}))
	resp, err := c.GetURL("http://localhost:8080", time.Second*5, WithUnmarshalResp(&data), WithByteResp(), WithHTTPResp(), WithReaderResp(), WithStringResp())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
	fmt.Printf("%+v\n", data)
	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("%+v\n", string(b))
}
