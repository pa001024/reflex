package weibo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pa001024/reflex/util"
)

const (
	_TSINA_OAUTH_VERSION = "2.0"
	_TSINA_OAUTH_BASE    = "https://api.weibo.com/oauth2/"
)

func (this *SinaWeibo) Authorize() (authurl string) {
	return _TSINA_OAUTH_BASE + "authorize?" + (url.Values{
		"client_id":     {this.AppKey},
		"redirect_uri":  {this.CallbackUrl},
		"response_type": {"code"},
		"display":       {"client"},
	}).Encode()
}
func (this *SinaWeibo) AccessToken(code string) (token string) {
	res, err := http.PostForm(_TSINA_OAUTH_BASE+"access_token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {this.AppKey},
			"client_secret": {this.AppSecret},
			"code":          {code},
			"redirect_uri":  {this.CallbackUrl},
		})
	if err != nil {
		util.ERROR.Log("Fail to AccessToken:", err)
		return
	}

	defer res.Body.Close()
	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)
	if body["error"] != nil || body["access_token"] == nil {
		util.ERROR.Log("Fail to AccessToken(Remote):", body["error"])
		return
	}
	this.Token = body["access_token"].(string)
	i, _ := strconv.Atoi(body["expires_in"].(string))
	ex := time.Now().Add(time.Duration(i) * time.Second)
	this.ExpiresIn = ex
	return this.Token
}
