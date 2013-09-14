package webqq

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/pa001024/reflex/util"
)

const (
	_WEBQQ_APPID  = "1003903"
	_WEBQQ_TYPE   = "10"
	_PTLOGIN_URL  = "https://ssl.ptlogin2.qq.com/"
	_JS_VER       = "10037"     // [2013-8-27] 10041
	_LOGIN_ACTION = "2-10-5837" // [2013-8-27] 4-9-7230
)

// [1]检查前获取sig, 永久 [2013-8-27]
func (this *WebQQ) ptlogin_login_sig() (login_sig string, err error) {
	res, err := this.client.Get(this.sig_url)
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	util.Try(err)
	res.Body.Close()
	login_sig = regexp.MustCompile(`var g_login_sig=encodeURIComponent\("(.+?)"\);`).FindStringSubmatch(string(bs))[1]
	util.DEBUG.Logf("[ptlogin_login_sig] login_sig = %s ", login_sig)
	return
}

// [2]检查, 可重复
func (this *WebQQ) ptlogin_check() (codetoken, code string, err error) {
	res, err := this.getWithReferer(_PTLOGIN_URL+"check?"+(url.Values{
		"uin":     {this.id_str},
		"appid":   {_WEBQQ_APPID},
		"js_ver":  {_JS_VER},
		"js_type": {"0"}, "u1": {"http://web2.qq.com/loginproxy.html"},
		"r": {rand_r()}}).Encode(), this.sig_url)
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	util.Try(err)
	s := string(bs)
	if s[14] == '0' { //ptui_checkVC('0','!YIZ','\x00\x00\x00\x00\x04\x82\x43\x7e');
		util.DEBUG.Log("[ptlogin_check] check ok")
		code = s[18:22]
	} else if s[14] == '1' { //ptui_checkVC('1','f5961ec58f7057466b40d950d5817b173f2400a1bca95ddc','\x00\x00\x00\x00\x2d\x16\xa2\x80');
		util.DEBUG.Log("[ptlogin_check] check need getimage")
		codetoken = s[18:66]
		code, err = this.ptlogin_getimage(codetoken)
	} else {
		err = fmt.Errorf("[ptlogin_check]失败返回值: %s", s)
	}
	util.Try(err)
	util.DEBUG.Logf("[ptlogin_check] return OK [code] %s", code)
	return
}

// [3]单点登录, 可重复
func (this *WebQQ) ptlogin_login(code string) (urlStr string, err error) {
	res, err := this.getWithReferer(
		_PTLOGIN_URL+"login?"+(url.Values{
			"u":            {this.id_str},
			"p":            {this.genPwd(code)},
			"verifycode":   {code},
			"webqq_type":   {_WEBQQ_TYPE},
			"action":       {_LOGIN_ACTION},
			"js_ver":       {_JS_VER},
			"aid":          {_WEBQQ_APPID},
			"remember_uin": {"1"}, "login2qq": {"1"}, "u1": {"http://web2.qq.com/loginproxy.html?login2qq=1&webqq_type=10"}, "h": {"1"}, "ptredirect": {"0"}, "ptlang": {"2052"}, "daid": {"164"}, "from_ui": {"1"}, "pttype": {"1"}, "dumy": {""}, "fp": {"loginerroralert"}, "mibao_css": {"m_webqq"}, "t": {"1"}, "g": {"1"}, "js_type": {"0"},
			"login_sig": {this.login_sig}}).Encode(), this.sig_url)
	util.Try(err)
	defer res.Body.Close()
	raw := util.MustReadAll(res.Body)
	ss := strings.Split(string(raw), "'")
	// 出错
	if ss[1] != "0" {
		err = fmt.Errorf("[ptlogin_login] fail_login %s:%s", ss[1], ss[9])
		return
	}
	//ptuiCB('0','0','http://ptlogin4.web2.qq.com/check_sig?pttype=1&uin=2735284921&service=login&nodirect=0&ptsig=MPlx81vcwwhHDYZeAsCdaFoQg3nTXyy67sQAYCewxu0_&s_url=http%3a%2f%2fweb2.qq.com%2floginproxy.html%3flogin2qq%3d1%26webqq%5ftype%3d10&f_url=&ptlang=2052&ptredirect=100&aid=1003903&daid=164&j_later=0&low_login_hour=0&regmaster=0','0','登录成功！', '菊菊菊菊菊菊');
	urlStr = ss[5]
	util.DEBUG.Logf("[ptlogin_login] return %s", ss[9])
	return
}

// [4]用获取cookie [2013-8-27]
func (this *WebQQ) ptlogin_check_sig(urlStr string) (err error) {
	res, err := this.client.Get(urlStr)
	util.Try(err)
	res.Body.Close()
	util.DEBUG.Logf("[ptlogin_check_sig] check ok")
	return
}
