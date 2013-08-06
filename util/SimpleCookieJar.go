package util

import (
	"net/http"
	"net/url"
)

func NewSimpleCookieJar() *SimpleCookieJar {
	return &SimpleCookieJar{make([]*http.Cookie, 0)}
}

type SimpleCookieJar struct{ cookies []*http.Cookie }

func (this *SimpleCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
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

func (this *SimpleCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return this.cookies
}
