package webqq

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/pa001024/MoeWorker/util"
)

const (
	WEBQQ_APPID  = "1003903"
	WEBQQ_TYPE   = "10"
	PTLOGIN_URL  = "https://ssl.ptlogin2.qq.com/"
	JS_VER       = "10037"     // [2013-8-27] 10041
	LOGIN_ACTION = "2-10-5837" // [2013-8-27] 4-9-7230
)

// [1]检查前获取sig, 永久 [2013-8-27]
func (this *WebQQ) ptlogin_login_sig() (login_sig string, err error) {
	res, err := this.client.Get(this.SigUrl)
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	util.Try(err)
	res.Body.Close()
	login_sig = regexp.MustCompile(`var g_login_sig=encodeURIComponent\("(.+?)"\);`).FindStringSubmatch(string(bs))[1]
	return
}

// [2]检查, 可重复
func (this *WebQQ) ptlogin_check() (codetoken, code, pwd string, err error) {
	util.DEBUG.Log("[ptlogin_check] Start")
	res, err := this.getWithReferer(PTLOGIN_URL+"check?"+(url.Values{
		"uin":     {this.IdStr},
		"appid":   {WEBQQ_APPID},
		"js_ver":  {JS_VER},
		"js_type": {"0"}, "u1": {"http://web2.qq.com/loginproxy.html"},
		"r": {rand_r()}}).Encode(), this.SigUrl)
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	util.Try(err)
	s := string(bs)
	if s[14] == '0' { //ptui_checkVC('0','!YIZ','\x00\x00\x00\x00\x04\x82\x43\x7e');
		code = s[18:22]
	} else if s[14] == '1' { //ptui_checkVC('1','f5961ec58f7057466b40d950d5817b173f2400a1bca95ddc','\x00\x00\x00\x00\x2d\x16\xa2\x80');
		codetoken = s[18:66]
		code, err = this.ptlogin_getimage(codetoken)
	} else {
		err = fmt.Errorf("[ptlogin_check]失败返回值: %s", s)
	}
	util.Try(err)
	pwd = this.genPwd(code)
	util.DEBUG.Log("[ptlogin_check] Success genPwd ", pwd)
	return
}

// [3]单点登录, 可重复
func (this *WebQQ) ptlogin_login(code, pwd string) (urlStr, msg string, err error) {
	res, err := this.getWithReferer(
		PTLOGIN_URL+"login?"+(url.Values{
			"u":            {this.IdStr},
			"p":            {pwd},
			"verifycode":   {code},
			"webqq_type":   {WEBQQ_TYPE},
			"action":       {LOGIN_ACTION},
			"js_ver":       {JS_VER},
			"aid":          {WEBQQ_APPID},
			"remember_uin": {"1"}, "login2qq": {"1"}, "u1": {"http://web2.qq.com/loginproxy.html?login2qq=1&webqq_type=10"}, "h": {"1"}, "ptredirect": {"0"}, "ptlang": {"2052"}, "daid": {"164"}, "from_ui": {"1"}, "pttype": {"1"}, "dumy": {""}, "fp": {"loginerroralert"}, "mibao_css": {"m_webqq"}, "t": {"1"}, "g": {"1"}, "js_type": {"0"},
			"login_sig": {this.LoginSig}}).Encode(), this.SigUrl)
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	util.Try(err)
	res.Body.Close()
	ss := strings.Split(string(bs), "'")
	if ss[1] == "4" {
		err = fmt.Errorf("%s", ss[9])
		util.Try(err)
	} else if ss[1] != "0" {
		err = fmt.Errorf("[ptlogin_login]失败返回值: %s", string(bs))
		util.Try(err)
	}
	//ptuiCB('0','0','http://web.qq.com/loginproxy.html?login2qq=1&webqq_type=10','0','登录成功！', '菊菊菊菊菊菊');
	//ptuiCB('0','0','http://ptlogin4.web2.qq.com/check_sig?pttype=1&uin=2735284921&service=login&nodirect=0&ptsig=MPlx81vcwwhHDYZeAsCdaFoQg3nTXyy67sQAYCewxu0_&s_url=http%3a%2f%2fweb2.qq.com%2floginproxy.html%3flogin2qq%3d1%26webqq%5ftype%3d10&f_url=&ptlang=2052&ptredirect=100&aid=1003903&daid=164&j_later=0&low_login_hour=0&regmaster=0','0','登录成功！', '菊菊菊菊菊菊');

	urlStr = ss[5]
	msg = ss[9]
	if strings.HasPrefix(urlStr, "http://web2.qq.com/loginproxy") {
		util.WARN.Log("[ptlogin_login] fail_check_sig")
	}
	return
}

// [4]用获取cookie [2013-8-27]
func (this *WebQQ) ptlogin_check_sig(oldurl string) (err error) {
	res, err := this.client.Get(oldurl)
	util.Try(err)
	res.Body.Close()
	return
}
