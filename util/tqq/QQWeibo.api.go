package tqq

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/pa001024/reflex/util"
)

const (
	_TQQ_API_BASE = "https://open.t.qq.com/api/"
)

// 腾讯微博API
type QQWeibo struct {
	AppKey       string    `json:"client_id"`     // AppKey
	AppSecret    string    `json:"client_secret"` // AppSecret
	CallbackUrl  string    `json:"redirect_uri"`  // 验证URL
	Token        string    `json:"access_token"`  // OAuth2.0 验证码
	ExpiresIn    time.Time `json:"expires_in"`    // 失效时间
	RefreshToken string    `json:"refresh_token"` // OAuth2.0 疼迅的二次验证码
	OpenID       string    `json:"openid"`        // 平台ID
	// ClientIP     string    `json:"clientip"`   // 改为动态计算
}

func (this *QQWeibo) postStatus(api string, args *url.Values) (rst *QQWeiboStatus, err error) {
	// OAuth required
	args.Set("oauth_consumer_key", this.AppKey)
	args.Set("access_token", this.Token)
	args.Set("openid", this.OpenID)
	args.Set("clientip", util.GetIP())
	res, err := http.PostForm(_TQQ_API_BASE+"t/"+api, *args)
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

// 发送微博
func (this *QQWeibo) Update(status string) (rst *QQWeiboStatus, err error) {
	rst, err = this.postStatus("add", &url.Values{
		"format":  {"json"},
		"content": {status},
	})
	return
}

// 转发微博
func (this *QQWeibo) Repost(status string, oriId string) (rst *QQWeiboStatus, err error) {
	rst, err = this.postStatus("re_add", &url.Values{
		"format":  {"json"},
		"content": {status},
		"reid":    {oriId},
	})
	return
}

// 删除微博
func (this *QQWeibo) Destroy(oriId string) (rst *QQWeiboStatus, err error) {
	rst, err = this.postStatus("del", &url.Values{
		"id": {oriId},
	})
	return
}

// 上传图片地址
func (this *QQWeibo) UploadUrl(status string, urlText string) (rst *QQWeiboStatus, err error) {
	rst, err = this.postStatus("add_pic_url", &url.Values{
		"format":  {"json"},
		"content": {status},
		"pic_url": {urlText},
	})
	return
}

// 上传图片
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

	res, err := http.Post(_TQQ_API_BASE+"add_pic", formdata.FormDataContentType(), &bpic)
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
