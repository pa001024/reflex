package tqq

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pa001024/reflex/util"
)

const (
	_TQQ_OAUTH_VERSION = "2.a"
	_TQQ_OAUTH_BASE    = "https://open.t.qq.com/cgi-bin/oauth2/"
)

// 返回Authorize网页地址
func (this *QQWeibo) Authorize() (authurl string) {
	return _TQQ_OAUTH_BASE + "authorize?" + (url.Values{
		"client_id":     {this.AppKey},
		"redirect_uri":  {this.CallbackUrl},
		"response_type": {"code"},
	}).Encode()
}

// 验证Authorize后返回的code
func (this *QQWeibo) AccessToken(code string) (token string) {
	res, err := http.PostForm(_TQQ_OAUTH_BASE+"access_token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {this.AppKey},      // yourappkey
			"client_secret": {this.AppSecret},   // yourpppsecret
			"code":          {code},             // xxxxxxxxxxxxxx
			"redirect_uri":  {this.CallbackUrl}, // http://some/weibocb.php
		})
	if err != nil {
		util.ERROR.Err("Fail to AccessToken:", err)
		return
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		util.ERROR.Err("Fail to AccessToken:", err)
		return
	}
	body, _ := url.ParseQuery(string(b))
	if body.Get("error") != "" || body.Get("access_token") == "" {
		util.ERROR.Err("Fail to AccessToken(Remote):", body.Get("error"))
		return
	}
	this.Token = body.Get("access_token")
	this.OpenID = body.Get("openid")
	this.RefreshToken = body.Get("refresh_token")
	i, _ := strconv.Atoi(body.Get("expires_in"))
	ex := time.Now().Add(time.Duration(i) * time.Second)
	this.ExpiresIn = ex
	return this.Token
}
