package weiyun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pa001024/reflex/util"
)

const (
	_WEIYUN_SHARE_URL      = "http://share.weiyun.com/"
	_WEIYUN_VIEW_SHARE_URL = "http://web2.cgi.weiyun.com/wy_share_v2.fcg?cmd=view_share&g_tk=&callback=CALLBACK&data="
	_WEIYUN_STORE_URL      = "http://web.cgi.weiyun.com/share_dl.fcg"
)

type ShareInfo struct {
	Createtime string   `json:"createtime"`
	DirList    []string `json:"dir_list"`
	Dlskey     string   `json:"dlskey"`
	Downcnt    int32    `json:"downcnt"`
	FileList   []struct {
		FileId   string `json:"file_id"`
		FileName string `json:"file_name"`
		FileSize string `json:"file_size"`
	} `json:"file_list"`
	NoteInfo  string `json:"note_info"`
	PdirKey   string `json:"pdir_key"`
	ResType   int32  `json:"res_type"`
	ShareFlag int32  `json:"share_flag"`
	ShareKey  string `json:"share_key"`
	Sharename string `json:"sharename"`
	ShortUrl  string `json:"short_url"`
	ShowStore bool   `json:"show_store"`
	Storecnt  int32  `json:"storecnt"`
	Uin       int64  `json:"uin"`
	Viewcnt   int32  `json:"viewcnt"`
}

// 从url地址获取sharekey 支持三种格式: 79ff0e64bd9c88cbd53a0ae72c14401c http://url.cn/PVX6FN(或其他短链,需跳转到后者) http://share.weiyun.com/79ff0e64bd9c88cbd53a0ae72c14401c
func ParseShareKey(rawurl string) (sharekey string, err error) {
	if len(rawurl) > 4 && rawurl[0:4] != "http" {
		sharekey = rawurl
		return
	}
	u, err := url.Parse(rawurl)
	if err != nil || u == nil {
		return
	}
	if u.Host != "share.weiyun.com" {
		res, _err := http.Get(rawurl)
		if _err != nil {
			return "", _err
		}
		u = res.Request.URL
	}
	if u != nil && len(u.Path) == 33 {
		sharekey = u.Path[1:]
	} else {
		err = fmt.Errorf("invalid rawurl %v", rawurl)
	}
	return
}

/*
	[2014.2.6] 1.获取文件信息
	GET http://web2.cgi.weiyun.com/wy_share_v2.fcg
	cmd:view_share
	g_tk:
	callback:CALLBACK
	data:{"req_header":{"cmd":"view_share","main_v":11,"proto_ver":10006,"sub_v":1,"encrypt":0,"msg_seq":1,"source":30111,"appid":30111,"client_ip":"127.0.0.1","token":"4d3754f563ad04a56fece81bbcc83302"},"req_body":{"share_key":"79ff0e64bd9c88cbd53a0ae72c14401c"}}
	_:1391326126322

	RET:
	try{CALLBACK({"rsp_body":{"createtime":"2013-08-28 12:40:43","dir_list":[],"dlskey":"74c726bfa1e9978207b41e9b312fa5d8709a1c5c7099774009af37429e0c58fd40e9cb52d7b84819","downcnt":54,"file_list":[{"file_id":"bdbc4714-b322-437b-87b3-25c978fcfa37","file_name":"jd-gui.exe","file_size":"719872"}],"note_info":null,"pdir_key":"80a2162d77c92765438ca4ef1d170515","res_type":0,"share_flag":0,"share_key":"79ff0e64bd9c88cbd53a0ae72c14401c","sharename":"jd-gui.exe","short_url":"http://url.cn/PVX6FN","show_store":true,"storecnt":0,"uin":756458112,"viewcnt":87},"rsp_header":{"ret":0}})}catch(e){};

*/
func GetShareInfo(sharekey string) (shareInfo ShareInfo, err error) {
	res, err := http.Get(_WEIYUN_VIEW_SHARE_URL +
		url.QueryEscape(util.ToJson(
			"req_header", struct {
				Cmd      string `json:"cmd"`
				MainV    int32  `json:"main_v"`
				ProtoVer int32  `json:"proto_ver"`
				SubV     int32  `json:"sub_v"`
				Encrypt  int32  `json:"encrypt"`
				MsgSeq   int32  `json:"msg_seq"`
				Source   int32  `json:"source"`
				Appid    int32  `json:"appid"`
				ClientIp string `json:"client_ip"`
				Token    string `json:"token"`
			}{"view_share", 11, 10006, 1, 0, 1, 30111, 30111, "127.0.0.1", "4d3754f563ad04a56fece81bbcc83302"},
			"req_body", struct {
				ShareKey string `json:"share_key"`
			}{sharekey},
		)) + "&_=" + util.JsCurrentTime())
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}
	rsp := &struct {
		RspBody   ShareInfo `json:"rsp_body"`
		RspHeader struct {
			Ret int32 `json:"ret"`
		} `json:"rsp_header"`
	}{}
	err = json.Unmarshal(b[13:len(b)-13], rsp) // try{CALLBACK({...})}catch(e){};
	shareInfo = rsp.RspBody
	return
}

/*
	[2014.2.6] 2.从shareinfo获取文件下载地址
	POST http://web.cgi.weiyun.com/share_dl.fcg
	sharekey:79ff0e64bd9c88cbd53a0ae72c14401c
	uin:756458112
	skey:
	fid:bdbc4714-b322-437b-87b3-25c978fcfa37
	dir:
	pdir:80a2162d77c92765438ca4ef1d170515
	zn:jd-gui.exe
	os_info:windows
	browser:chrome
	ver:12
	err_callback:http://www.weiyun.com/web/callback/iframe_share_down_fail.html

	RET 302 -> true file download url
*/
func GetShareFile(shareInfo ShareInfo, index int) (dlurl string, err error) {
	if index >= len(shareInfo.FileList) {
		err = fmt.Errorf("index out of files")
		return
	}
	req, err := http.NewRequest("POST", _WEIYUN_STORE_URL, strings.NewReader((url.Values{
		"sharekey":     {shareInfo.ShareKey},
		"uin":          {""},
		"skey":         {""},
		"fid":          {shareInfo.FileList[index].FileId},
		"dir":          {""},
		"pdir":         {shareInfo.PdirKey},
		"zn":           {shareInfo.FileList[index].FileId},
		"os_info":      {"windows"},
		"browser":      {"chrome"},
		"ver":          {"12"},
		"err_callback": {"http://www.weiyun.com/web/callback/iframe_share_down_fail.htm"},
	}).Encode()))
	req.Header.Add("Referer", _WEIYUN_SHARE_URL+shareInfo.ShareKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	res.Body.Close()
	dlurl = res.Request.URL.String() // 302
	return
}

// 一步到位
func ParseAndDL(raw string) (dl string, err error) {
	sharekey, err := ParseShareKey(raw)
	if err != nil {
		return
	}
	var info ShareInfo
	info, err = GetShareInfo(sharekey)
	if err != nil {
		return
	}
	if len(info.FileList) > 0 {
		dl, err = GetShareFile(info, 0)
	}
	return
}
