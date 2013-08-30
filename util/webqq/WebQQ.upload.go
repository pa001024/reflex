package webqq

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"regexp"

	"github.com/pa001024/MoeWorker/util"
)

const (
	_CFACE_UPLOAD_URL   = "http://up.web2.qq.com/cgi-bin/cface_upload"
	_OFFLINE_UPLOAD_URL = "http://weboffline.ftn.qq.com/ftn_access/upload_offline_pic"
)

var (
	// <head><script type="text/javascript">document.domain='qq.com';parent.EQQ.Model.ChatMsg.callbackSendPicGroup({'ret':0,'msg':'C5868D1266F58FA657A9352C76477296.jPgC5868D1266F58FA657A9352C76477296.jPg -7001'});</script></head><body></body>
	cfacePatten = regexp.MustCompile(`EQQ\.Model\.ChatMsg\.callbackSendPicGroup\(\{'ret':0,'msg':'(.+?\..{3})`)
	// <head><script type="text/javascript">document.domain='qq.com';parent.EQQ.Model.ChatMsg.callbackSendPic({"retcode":0, "result":"OK", "progress":100, "filesize":2730, "fileid":"3", "filename":"50.gif", "filepath":"/c501633b-4f26-4af3-a714-2f3db4134500"});</script></head><body></body>
	offlinePatten = regexp.MustCompile(`EQQ\.Model\.ChatMsg\.callbackSendPicGroup\((.+)\);`)
)

// 上传自定义表情
func (this *WebQQ) UploadCustomFace(pic io.Reader) (v string, err error) {
	form := &bytes.Buffer{}
	formdata := multipart.NewWriter(form)
	formdata.WriteField("from", "control")
	formdata.WriteField("f", "EQQ.Model.ChatMsg.callbackSendPicGroup")
	formdata.WriteField("vfwebqq", this.vfwebqq)
	formdata.WriteField("fileid", util.ToString(fileid))
	fileid++
	picdata, _ := formdata.CreateFormFile("custom_face", "1.png")
	io.Copy(picdata, pic)
	formdata.Close()
	res, err := this.client.Post(_CFACE_UPLOAD_URL+"?time="+util.JsCurrentTime(), formdata.FormDataContentType(), form)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	sb := cfacePatten.FindStringSubmatch(string(b))
	if len(sb) > 1 {
		v = sb[1]
	}
	return
}

var (
	fileid = 1
)

// 上传离线图片
func (this *WebQQ) UploadOfflinePic(peeruin Uin, pic io.Reader) (v *ResultOfflinePic, err error) {
	form := &bytes.Buffer{}
	formdata := multipart.NewWriter(form)
	formdata.WriteField("callback", "parent.EQQ.Model.ChatMsg.callbackSendPic")
	formdata.WriteField("locallangid", "2052")
	formdata.WriteField("clientversion", "1409")
	formdata.WriteField("uin", this.id_str)
	formdata.WriteField("skey", this.getCookie(util.MustParseUrl(_OFFLINE_UPLOAD_URL), "skey"))
	formdata.WriteField("appid", _WEBQQ_APPID)
	formdata.WriteField("peeruin", peeruin.String())
	formdata.WriteField("vfwebqq", this.vfwebqq)
	formdata.WriteField("fileid", util.ToString(fileid))
	formdata.WriteField("senderviplevel", "0")
	formdata.WriteField("reciverviplevel", "0")
	fileid++
	picdata, _ := formdata.CreateFormFile("file", "1.png")
	io.Copy(picdata, pic)
	formdata.Close()
	res, err := this.client.Post(_CFACE_UPLOAD_URL+"?time="+util.JsCurrentTime(), formdata.FormDataContentType(), form)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	sb := offlinePatten.FindStringSubmatch(string(b))
	if len(sb) > 1 {
		v = &ResultOfflinePic{}
		json.Unmarshal([]byte(sb[1]), v)
	}
	return
}

// 上传离线图片 result结构
type ResultOfflinePic struct {
	Code     int32  `json:"retcode"`
	Result   string `json:"result"`
	Progress int32  `json:"progress"`
	Filesize int32  `json:"filesize"`
	Filename int32  `json:"filename"`
	Filepath int32  `json:"filepath"`
}
