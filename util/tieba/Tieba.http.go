package tieba

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// http util
func (this *Tieba) post(urlStr string, w io.Reader, cookie string) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", urlStr, w)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("User-Agent", "BaiduTieba for Android 5.1.3")
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}
	ch := make(chan bool)
	go func() {
		res, err = this.client.Do(req)
		if res != nil {
			this.client.Jar.SetCookies(req.URL, res.Cookies())
		}
		<-ch
	}()
	select {
	case ch <- true:
	case <-time.After(60 * time.Second):
		http.DefaultTransport.(*http.Transport).CancelRequest(req)
		err = fmt.Errorf("Login Timeout")
	}
	close(ch)
	return
}
