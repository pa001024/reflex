package tieba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pa001024/reflex/util"
)

// 获取贴吧fid
func (this *Tieba) GetFid(kw string) (fid string) {
	res, err := this.client.Get(_FID_URL + kw)
	if err != nil {
		return
	}
	defer res.Body.Close()
	str := string(util.MustReadAll(res.Body))
	rex := regexp.MustCompile(`"fid" +?value="(\d+)"`)
	m := rex.FindStringSubmatch(str)
	if len(m) > 0 {
		fid = m[1]
	}
	return
}

// [未完成]发帖
func (this *Tieba) postFid(fid, kw, title, content string) (result *Result, err error) {
	this.checkTbs()

	parm := []string{
		this._BDUSS(),
		this._CLIENT_ID(),
		_CLIENT_TYPE,
		_CLIENT_VERSION,
		_PHONE_IMEI,
		"content=" + content,
		"fid=" + fid,
		"kw=" + kw,
		this._TBS(),
		"stErrorNums=0",
		"timestamp=1390753277530",
		"title=" + title,
	}
	str := fmt.Sprint(strings.Join(parm, "&"),
		"&sign=", this.getSign(parm))
	req := bytes.NewBufferString(str)
	res, err := this.post(_THREAD_URL, req, "BDUSS="+this.BDUSS+";")
	bin := util.MustReadAll(res.Body)
	util.DEBUG.Log("[Post]", str, "\n========\n", string(bin))
	result = &Result{}
	err = json.Unmarshal(bin, result)
	return
}

// [未完成]发帖
func (this *Tieba) _Post(kw, title, content string) (result *Result, err error) {
	fid := this.GetFid(kw)
	return this.postFid(fid, kw, title, content)
}

// [未完成]回帖
func (this *Tieba) _Reply(qid, kw, tid, content string) (result *Result, err error) {
	this.checkTbs()
	parm := []string{
		this._BDUSS(),
		"kw=" + kw,
		this._TBS(),
	}
	str := fmt.Sprint(strings.Join(parm, "&"),
		"tid=", tid,
		"content=", content,
		"&sign=", this.getSign(parm))
	req := bytes.NewBufferString(str)
	res, err := this.post(_POST_URL, req, "")
	bin := util.MustReadAll(res.Body)
	util.DEBUG.Log("[Reply]", str, "\n========\n", string(bin))
	result = &Result{}
	err = json.Unmarshal(bin, result)
	return
}
