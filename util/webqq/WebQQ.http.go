package webqq

import (
	"bytes"
	"net/http"
	"net/url"
)

// 获取特定cookie值
func (this *WebQQ) GetCookie(url *url.URL, name string) (ret string) {
	if this.client.Jar != nil {
		for _, v := range this.client.Jar.Cookies(url) {
			DEBUG.Log(v)
			if v.Name == name {
				ret = v.Value
				// return
			}
		}
	}
	return
}

// 带变量Referer GET
func (this *WebQQ) GetWithReferer(urlStr string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if this.client.Jar != nil {
		for _, v := range this.client.Jar.Cookies(req.URL) {
			req.AddCookie(v)
		}
	}
	req.Header.Add("Referer", "http://s.web2.qq.com/proxy.html?v=20110412001&callback=1&id=3")
	return this.client.Do(req)
}

// 带固定Referer POST
func (this *WebQQ) PostFormWithReferer(url string, val url.Values) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(val.Encode()))
	if this.client.Jar != nil {
		for _, v := range this.client.Jar.Cookies(req.URL) {
			req.AddCookie(v)
		}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://s.web2.qq.com/proxy.html?v=20110412001&callback=1&id=1")
	return this.client.Do(req)
}
