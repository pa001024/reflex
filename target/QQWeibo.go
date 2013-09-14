package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pa001024/reflex/source"
	"github.com/pa001024/reflex/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	TQQ_OAUTH_VERSION = "2.a"
)

func (this *QQWeibo) Send(src *source.FeedInfo) (rid string, e error) {
	util.DEBUG.Logf("QQWeibo.Send(%v:%v,repost_id:%v,title:%v,content:%v,picurl:%v)", src.SourceId, src.Id, src.RepostId, src.Title, src.Content, src.PicUrl)
	if src.RepostId != "" {
		r, err := this.Repost(src.Content, src.RepostId)
		if err != nil {
			e = err
			return
		}
		util.INFO.Logf("[qq.%v] Repost sent: %v", this.Name, r.Url())
		return util.ToString(r.Id), nil
	} else if src.PicUrl != nil && len(src.PicUrl) > 0 {
		r, err := this.UploadUrl(src.Content, src.PicUrl[0])
		if err != nil {
			e = err
			return
		}
		util.INFO.Logf("[qq.%v] UploadUrl sent: %v", this.Name, r.Url())
		return util.ToString(r.Id), nil
	} else {
		r, err := this.Update(src.Content)
		if err != nil {
			e = err
			return
		}
		util.INFO.Logf("[qq.%v] Update sent: %v", this.Name, r.Url())
		return util.ToString(r.Id), nil
	}
	return
}

func (this *QQWeibo) GetMethod() []*TargetMethod { return this.Method }
func (this *QQWeibo) GetId() string              { return this.Name }

type QQWeibo struct { // 腾讯微博API 实现接口IWeibo
	IWeibo
	ITarget
	Target

	AppKey       string    `json:"client_id"`     // AppKey
	AppSecret    string    `json:"client_secret"` // AppSecret
	CallbackUrl  string    `json:"redirect_uri"`  // 验证URL
	Token        string    `json:"access_token"`  // OAuth2.0 验证码
	ExpiresIn    time.Time `json:"expires_in"`    // 失效时间
	RefreshToken string    `json:"refresh_token"` // OAuth2.0 疼迅的二次验证码
	OpenID       string    `json:"openid"`        // 平台ID
	// ClientIP     string    `json:"clientip"`   // 改为动态计算
}
type QQWeiboResult struct {
	ErrorCode int            `json:"errcode"` // 错误代码
	Error     string         `json:"msg"`     // 返回信息
	Ret       int            `json:"ret"`     // 返回值
	Data      *QQWeiboStatus `json:"data"`    // 数据
	// SeqId     string         `json:"seqid"`   // 序列号 (无需使用)
}
type QQWeiboStatus struct {
	IStatus
	Id        int64  `json:"id"`        // 微博id
	CreatedAt string `json:"timestamp"` // 微博发表时间
}

func (this *QQWeibo) Authorize() (authurl string) {
	return "https://open.t.qq.com/cgi-bin/oauth2/authorize?" + (url.Values{
		"client_id":     {this.AppKey},
		"redirect_uri":  {this.CallbackUrl},
		"response_type": {"code"},
	}).Encode()
}
func (this *QQWeibo) AccessToken(code string) (token string) {
	res, err := http.PostForm("https://open.t.qq.com/cgi-bin/oauth2/access_token",
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
func (this *QQWeibo) PostStatus(api string, args *url.Values) (rst *QQWeiboStatus, err error) {
	// OAuth
	args.Set("oauth_consumer_key", this.AppKey)
	args.Set("access_token", this.Token)
	args.Set("openid", this.OpenID)
	args.Set("clientip", util.GetIP())
	res, err := http.PostForm("https://open.t.qq.com/api/t/"+api, *args)
	if err != nil {
		util.ERROR.Errf("Error on call %s: %v", api, err)
		return
	}
	defer res.Body.Close()
	dst := &QQWeiboResult{}
	json.NewDecoder(res.Body).Decode(dst)
	if dst.Error != "" {
		util.ERROR.Errf("Error on call %s(Remote): %v : %v", api, dst.ErrorCode, dst.Error)
		return nil, RemoteError(dst.Error)
	}
	rst = dst.Data
	return
}
func (this *QQWeibo) Update(status string) (rst *QQWeiboStatus, err error) {
	rst, err = this.PostStatus("add", &url.Values{
		"format":  {"json"},
		"content": {status},
	})
	return
}
func (this *QQWeibo) Repost(status string, oriId string) (rst *QQWeiboStatus, err error) {
	rst, err = this.PostStatus("re_add", &url.Values{
		"format":  {"json"},
		"content": {status},
		"reid":    {oriId},
	})
	return
}
func (this *QQWeibo) Destroy(oriId string) (rst *QQWeiboStatus, err error) {
	rst, err = this.PostStatus("del", &url.Values{
		"id": {oriId},
	})
	return
}
func (this *QQWeibo) UploadUrl(status string, urlText string) (rst *QQWeiboStatus, err error) {
	rst, err = this.PostStatus("add_pic_url", &url.Values{
		"format":  {"json"},
		"content": {status},
		"pic_url": {urlText},
	})
	return
}
func (this *QQWeibo) Upload(status string, pic io.Reader) (rst *QQWeiboStatus, err error) {
	// multipart/form-data
	var bpic bytes.Buffer
	formdata := multipart.NewWriter(&bpic)
	formdata.WriteField("oauth_consumer_key", this.AppKey)
	formdata.WriteField("access_token", this.Token)
	formdata.WriteField("openid", this.OpenID)
	formdata.WriteField("clientip", util.GetIP())

	formdata.WriteField("format", "json")
	formdata.WriteField("content", status)
	picdata, _ := formdata.CreateFormFile("pic", "image.png")
	io.Copy(picdata, pic)
	formdata.Close()

	res, err := http.Post("https://open.t.qq.com/api/t/add_pic", formdata.FormDataContentType(), &bpic)
	if err != nil {
		util.ERROR.Err("Error on call upload :", err)
		return
	}
	defer res.Body.Close()
	dst := &QQWeiboResult{}
	json.NewDecoder(res.Body).Decode(dst)
	if dst.Error != "" {
		util.ERROR.Err("Error on call upload (Remote):", dst.ErrorCode, ":", dst.Error)
		return nil, RemoteError(dst.Error)
	}
	rst = dst.Data
	return
}

func (this *QQWeiboStatus) Url() (urlText string) {
	urlText = fmt.Sprintf("http://t.qq.com/p/t/%v", this.Id)
	return
}
