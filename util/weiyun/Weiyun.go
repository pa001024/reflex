package weiyun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const (
	_WEIYUN_SHARE_URL = "http://share.weiyun.com/"
	_WEIYUN_STORE_URL = "http://sync.box.qq.com/share_dl.fcg?skey=&dir=&os_info=windows&browser=chrome&ver=12&"
)

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

// 从sharekey获取文件下载地址
func GetShareFile(sharekey string) (dlurl string, err error) {
	originurl := _WEIYUN_SHARE_URL + sharekey
	res, err := http.Get(originurl)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()
	ex := regexp.MustCompile(`var shareInfo = (.+);`).FindStringSubmatch(string(b))
	if len(ex) < 2 {
		err = fmt.Errorf("Error fetch shareurl with sharekey: %s", sharekey)
		return
	}
	type ShareInfo struct {
		Filename string `json:"filename"`
		FileList []struct {
			FileId string `json:"file_id"`
			// file_name
			// file_size
		} `json:"file_list"`
		PdirKey string `json:"pdir_key"`
		Uin     uint64 `json:"uin"`
	}
	shareInfo := &ShareInfo{}
	json.Unmarshal([]byte(ex[1]), shareInfo)
	storeurl := _WEIYUN_STORE_URL + (url.Values{
		"sharekey": {sharekey},
		"uin":      {fmt.Sprint(shareInfo.Uin)},
		"fid":      {shareInfo.FileList[0].FileId},
		"pdir":     {shareInfo.PdirKey},
		"zn":       {shareInfo.Filename},
	}).Encode()
	req, err := http.NewRequest("GET", storeurl, nil)
	req.Header.Add("Referer", originurl)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	res.Body.Close()
	dlurl = res.Request.URL.String()
	return
}

// 一步到位
func ParseAndDL(raw string) (dl string, err error) {
	sharekey, err := ParseShareKey(raw)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	dl, err = GetShareFile(sharekey)
	return
}
