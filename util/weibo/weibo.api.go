package weibo

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pa001024/reflex/util"
)

const (
	_TSINA_API_BASE = "https://api.weibo.com/2/"
)

func (this *SinaWeibo) postStatus(api string, args *url.Values) (rst *SinaWeiboStatus, err error) {
	args.Set("access_token", this.Token)
	res, err := http.PostForm(_TSINA_API_BASE+"statuses/"+api+".json", *args)
	if err != nil {
		util.ERROR.Log("Error on call", api+":", err)
		return
	}
	defer res.Body.Close()
	rst = &SinaWeiboStatus{}
	json.NewDecoder(res.Body).Decode(rst)
	if rst.Error != "" {
		return nil, RemoteError(rst.Error)
	}
	return
}
func (this *SinaWeibo) Update(status string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.postStatus("update", &url.Values{
		"status": {status},
	})
	return
}
func (this *SinaWeibo) Repost(status string, oriId string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.postStatus("repost", &url.Values{
		"status": {status},
		"id":     {oriId},
	})
	return
}
func (this *SinaWeibo) Destroy(oriId string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.postStatus("destroy", &url.Values{
		"id": {oriId},
	})
	return
}
func (this *SinaWeibo) Upload(status string, pic io.Reader) (rst *SinaWeiboStatus, err error) {
	// multipart/form-data
	var bpic bytes.Buffer
	formdata := multipart.NewWriter(&bpic)
	formdata.WriteField("access_token", this.Token)
	formdata.WriteField("status", status)
	picdata, _ := formdata.CreateFormFile("pic", "image.png")
	io.Copy(picdata, pic)
	formdata.Close()

	res, err := http.Post(_TSINA_API_BASE+"statuses/upload.json", formdata.FormDataContentType(), &bpic)
	if err != nil {
		util.ERROR.Log("Error on call upload :", err)
		return
	}
	defer res.Body.Close()
	rst = &SinaWeiboStatus{}
	json.NewDecoder(res.Body).Decode(rst)
	if rst.Error != "" {
		return nil, RemoteError(rst.Error)
	}
	return
}
func (this *SinaWeibo) UploadUrl(status string, urlText string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.postStatus("upload_url_text", &url.Values{
		"status": {status},
		"url":    {urlText},
	})
	return
}
func (this *SinaWeibo) ShortUrl(urlLong string) (rst string) {
	res, err := http.Get(_TSINA_API_BASE + "short_url/shorten.json?" + (url.Values{
		"access_token": {this.Token},
		"url_long":     {urlLong},
	}).Encode())
	if err != nil {
		util.ERROR.Log("Error on call ShortUrl:", err)
		return
	}
	defer res.Body.Close()
	v := &SinaWeiboShortUrlResult{}
	json.NewDecoder(res.Body).Decode(v)
	if v.Error != "" {
		util.ERROR.Log("Error on call ShortUrl[Remote]:", v.Error)
		return urlLong
	}
	return
}
