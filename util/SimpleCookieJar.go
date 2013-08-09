package util

import (
	"net/http"
	"net/url"
)

// 创建新的简单Jar
func NewSimpleCookieJar() *SimpleCookieJar {
	return &SimpleCookieJar{make([]*http.Cookie, 0)}
}

// 无域简单Jar
type SimpleCookieJar struct{ cookies []*http.Cookie }

// 实现接口 http.CookieJar.SetCookies(u, cookies)
func (this *SimpleCookieJar) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	found := false
	for a, s := range cookies {
		for a2, s2 := range this.cookies {
			if s.Name == s2.Name {
				this.cookies[a2] = cookies[a]
				found = true
				break
			}
		}
		if found == false {
			this.cookies = append(this.cookies, s)
		} else {
			found = false
		}
	}
}

// 实现接口 http.CookieJar.Cookies(u)
func (this *SimpleCookieJar) Cookies(_ *url.URL) []*http.Cookie {
	return this.cookies
}
