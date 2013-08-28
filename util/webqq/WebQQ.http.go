package webqq

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	// "github.com/pa001024/MoeWorker/util"
)

// 获取特定cookie值
func (this *WebQQ) getCookie(url *url.URL, name string) (ret string) {
	for _, v := range this.client.Jar.Cookies(url) {
		if v.Name == name {
			return v.Value
		}
	}
	return
}

// 带参数Referer GET
func (this *WebQQ) getWithReferer(urlStr, referer string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	for _, v := range this.client.Jar.Cookies(req.URL) {
		req.AddCookie(v)
	}
	req.Header.Add("Referer", referer)
	ch := make(chan bool)
	go func() {
		res, err = this.client.Do(req)
		this.client.Jar.SetCookies(req.URL, res.Cookies())
		<-ch
	}()
	select {
	case ch <- true:
		http.DefaultTransport.(*http.Transport).CancelRequest(req)
	case <-time.After(60 * time.Second):
		err = fmt.Errorf("Timeout POST %s", urlStr)
	}
	close(ch)
	return
}

// 带参数Referer POST
func (this *WebQQ) postFormWithReferer(urlStr, referer string, val url.Values) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(val.Encode()))
	for _, v := range this.client.Jar.Cookies(req.URL) {
		req.AddCookie(v)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", referer)
	ch := make(chan bool)
	go func() {
		res, err = this.client.Do(req)
		this.client.Jar.SetCookies(req.URL, res.Cookies())
		<-ch
	}()
	select {
	case ch <- true:
		http.DefaultTransport.(*http.Transport).CancelRequest(req)
	case <-time.After(60 * time.Second):
		err = fmt.Errorf("Timeout POST %s", urlStr)
	}
	close(ch)
	return
}
