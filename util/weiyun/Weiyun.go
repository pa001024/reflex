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
	_WEIYUN_DL_URL    = "http://sync.box.qq.com/share_dl.fcg?skey=&dir=&os_info=windows&browser=chrome&ver=12&"
)

type ShareInfo struct {
	Filename string `json:"filename"`
	FileList []struct {
		FileId string `json:"file_id"`
		// file_name
		// file_size
	} `json:"file_list"`
	PdirKey string `json:"pdir_key"`
	Uin     string `json:"uin"`
}

func HandleFile(sharekey string) (trueurl string, err error) {
	res, err := http.Get(_WEIYUN_SHARE_URL + sharekey)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	ex := regexp.MustCompile(`var shareInfo = (.+);`).FindStringSubmatch(string(b))
	if len(ex) < 2 {
		err = fmt.Errorf("Error fetch share url sharekey: %s", sharekey)
		return
	}
	var shareInfo *ShareInfo
	json.Unmarshal([]byte(ex[1]), sharekey)
	trueurl = _WEIYUN_DL_URL + (url.Values{
		"sharekey": {sharekey},
		"uin":      {shareInfo.Uin},
		"fid":      {shareInfo.FileList[0].FileId},
		"pdir":     {shareInfo.PdirKey},
		"zn":       {shareInfo.Filename},
	}).Encode()
	return
}
