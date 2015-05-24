package tieba

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"strings"

	"github.com/pa001024/reflex/util"
)

func (this *Tieba) login(vcode, vcode_md5 string) (result *LoginResult, err error) {
	parm := []string{
		"_client_id=" + this.client_id,
		_CLIENT_TYPE,
		_CLIENT_VERSION,
		_PHONE_IMEI,
		"passwd=" + base64.StdEncoding.EncodeToString([]byte(this.Password)),
		"un=" + url.QueryEscape(this.Username),
		"vcode_md5=" + vcode_md5,
	}
	code := ""
	if vcode != "" {
		code = "&vcode=" + vcode
	}
	str := fmt.Sprint(strings.Join(parm, "&"), code, "&sign=", this.getSign(parm))
	req := bytes.NewBufferString(str)
	res, err := this.post(_LOGIN_URL, req, "")
	bin := util.MustReadAll(res.Body)
	util.DEBUG.Log(str, "\n========\n", string(bin))
	result = &LoginResult{}
	err = json.Unmarshal(bin, result)
	return
}

func (this *Tieba) Login() (err error) {
	result, err2 := this.login("", "")
	err = err2
reset:
	if result.ErrorCode == 5 && result.Anti.NeedVcode == 1 {
		util.INFO.Log("[首次登录尝试失败]", result.ErrorMsg)
		vcode := ""
	retype:
		img, _, _ := image.Decode(util.FetchImageAsStream(result.Anti.VcodePicUrl))
		vcode = this.VcodeHandler(img)
		if vcode == "" {
			goto retype
		}
		result, err = this.login(vcode, result.Anti.VcodeMd5)
	}
	if result.ErrorCode != 0 {
		if result.ErrorCode != 5 {
			util.INFO.Log(result.ErrorMsg)
		}
		goto reset
	}
	this.BDUSS = result.User.BDUSS
	return
}
