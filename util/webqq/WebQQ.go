package webqq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/pa001024/MoeWorker/util"
	asc "github.com/pa001024/MoeWorker/util/ascgen"
)

var (
	DEBUG = util.NewLogger(false, "[DEBUG]")
	WARN  = util.NewLogger(true, "[WARN] ")
	INFO  = util.NewLogger(true, "[INFO] ")
)

type WebQQ struct {
	// ITarget
	// Target
	C *http.Client

	Id        string
	PasswdMd5 string
	// 瞬态
	ClientId   string
	VerifyCode string
	SessionId  string
	LoginSig   string
}

const (
	WEBQQ_APPID = "1003903"
	WEBQQ_TYPE  = "10"
	PTLOGIN_URL = "https://ssl.ptlogin2.qq.com/"
	CAPTCHA_URL = "http://captcha.qq.com/"
	CHANNEL_URL = "http://d.web2.qq.com/channel/"
)

func NewWebQQ(uin, pwd string) (this *WebQQ) {
	jar, _ := cookiejar.New(nil)
	this = &WebQQ{&http.Client{nil, nil, jar}, uin, pwd, fmt.Sprint(rand.Int31n(90000000) + 10000000), "", "", ""}
	return
}

// [1]检查前获取sig, 永久
func (this *WebQQ) ptlogin_login_sig() (login_sig string, err error) {
	res, err := this.C.Get("https://ui.ptlogin2.qq.com/cgi-bin/login?daid=164&target=self&style=5&mibao_css=m_webqq&appid=1003903&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fweb2.qq.com%2Floginproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=2013072300")
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	util.Try(err)
	res.Body.Close()
	login_sig = regexp.MustCompile(`var g_login_sig=encodeURIComponent\("(.+?)"\);`).FindStringSubmatch(string(bs))[1]
	return
}

// [2]检查, 可重复
func (this *WebQQ) ptlogin_check() (codetoken, code, pwd string, err error) {
	DEBUG.Log("[ptlogin_check] Start")
	res, err := this.C.Get(PTLOGIN_URL + "check?" + (url.Values{"uin": {this.Id},
		"appid":   {WEBQQ_APPID},
		"js_ver":  {"10037"},
		"js_type": {"0"},
		"u1":      {"http://web2.qq.com/loginproxy.html"},
		"r":       {fmt.Sprint(rand.ExpFloat64())}}).Encode())
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}
	s := string(bs)
	salt := ""
	if s[14] == '0' { //ptui_checkVC('0','!YIZ','\x00\x00\x00\x00\x04\x82\x43\x7e');
		code, salt = s[18:22], util.DecodeJsHex(s[25:len(s)-3])
	} else if s[14] == '1' { //ptui_checkVC('1','f5961ec58f7057466b40d950d5817b173f2400a1bca95ddc','\x00\x00\x00\x00\x2d\x16\xa2\x80');
		codetoken = s[18:66]
		code, err = this.ptlogin_getimage(codetoken)
		salt = util.DecodeJsHex(s[69 : len(s)-3])
	} else {
		err = errors.New("[ptlogin_check]失败返回值:" + s)
	}
	if err != nil {
		return
	}
	pwd = this.GenPwd(salt, code)
	DEBUG.Log("[ptlogin_check] Success GenPwd", pwd)
	return
}

// [2.1]获取验证码, 可重复
func (this *WebQQ) ptlogin_getimage(vCode string) (code string, err error) {
	res, err := this.C.Get(CAPTCHA_URL + "getimage?" + (url.Values{"uin": {this.Id}, "aid": {WEBQQ_APPID}, "r": {fmt.Sprint(rand.ExpFloat64())}, "vc_type": {vCode}}).Encode())
	if err != nil {
		return
	}
	defer res.Body.Close()
	rw, err := os.Create("vc.jpg")
	if err != nil {
		return
	}
	io.Copy(rw, res.Body)
	rw.Seek(0, 0)
	asc.ShowFile(os.Stdout, rw, asc.Console{6, 14, 120}, true, false)
	rw.Close()
	fmt.Print("Enter Verify Code: ")
	fmt.Scanf("%s", &code)
	return
}

// [3]单点登录, 可重复
func (this *WebQQ) ptlogin_login(code, pwd string) (err error) {
	res, err := this.GetWithReferer(
		PTLOGIN_URL+"login?"+(url.Values{"u": {this.Id},
			"p":            {pwd},
			"verifycode":   {code},
			"webqq_type":   {WEBQQ_TYPE},
			"remember_uin": {"1"},
			"login2qq":     {"1"},
			"aid":          {WEBQQ_APPID},
			"u1":           {"http://web2.qq.com/loginproxy.html?login2qq=1&webqq_type=10"},
			"h":            {"1"},
			"ptredirect":   {"0"},
			"ptlang":       {"2052"},
			"daid":         {"164"},
			"from_ui":      {"1"},
			"pttype":       {"1"},
			"dumy":         {""},
			"fp":           {"loginerroralert"},
			"action":       {"4-15-8246"},
			"mibao_css":    {"m_webqq"},
			"t":            {"1"},
			"g":            {"1"},
			"js_type":      {"0"},
			"js_ver":       {"10037"}}).Encode(),
		"https://ui.ptlogin2.qq.com/cgi-bin/login?daid=164&target=self&style=5&mibao_css=m_webqq&appid=1003903&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fweb2.qq.com%2Floginproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=20130723001")
	util.Try(err)
	bs, err := ioutil.ReadAll(res.Body)
	util.Try(err)
	res.Body.Close()
	ss := strings.Split(string(bs), "'")
	if ss[1] == "4" {
		return errors.New(ss[9])
	} else if ss[1] != "0" {
		err = errors.New("[ptlogin_login]失败返回值:" + string(bs))
		util.Try(err)
	}
	//ptuiCB('0','0','http://web.qq.com/loginproxy.html?login2qq=1&webqq_type=10','0','登录成功！', '菊菊菊菊菊菊');
	//ptuiCB('0','0','https://ssl.ptlogin2.qq.com/pt4_302?u1=http%3A//ptlogin4.web2.qq.com/check_sig%3Fpttype%3D1%26uin%3D2735284921%26service%3Dlogin%26nodirect%3D0%26ptsig%3DjHDCjZ5Our13vq2Kmx8VjeKxbVqg*UjyI01f2oGT8MY_%26s_url%3Dhttp%253a%252f%252fweb2.qq.com%252floginproxy.html%253flogin2qq%253d1%2526webqq%255ftype%253d10%26f_url%3D%26ptlang%3D2052%26ptredirect%3D100%26aid%3D1003903%26daid%3D164%26j_later%3D0%26low_login_hour%3D0%26regmaster%3D0','0','登录成功！', '菊菊菊菊菊菊');
	url := ss[5]
	msg := ss[9]
	if url == "http://web2.qq.com/loginproxy.html?login2qq=1&webqq_type=10" {
		WARN.Log("[ptlogin_login] fail_pt4_302")
		return
	} else {
		DEBUG.Log("[ptlogin_login]", msg)
		url, err = this.ptlogin_pt4_302(url)
		util.Try(err)
	}
	return
}

// [4]获取网址
func (this *WebQQ) ptlogin_pt4_302(oldurl string) (url string, err error) {
	pt4_302 := util.MustParseUrl(oldurl)
	check_sig := pt4_302.Query().Get("u1")
	if check_sig != "" {
		DEBUG.Log("[check_sig]", check_sig)
		return this.ptlogin_check_sig(check_sig)
	}
	return
}

// [5]获取cookie
func (this *WebQQ) ptlogin_check_sig(oldurl string) (url string, err error) {
	res, err := this.C.Get(oldurl)
	if err != nil {
		return
	}
	res.Body.Close()
	return
}

// [6]用令牌登录WebQQ
func (this *WebQQ) channel_login2(ptwebqq string) (hr *Login2Result, err error) {
	res, err := this.PostFormWithReferer(CHANNEL_URL+"login2",
		url.Values{
			"r": {util.ToJson(
				"status", "online",
				"ptwebqq", ptwebqq,
				"passwd_sig", "",
				"clientid", this.ClientId,
				"psessionid", nil,
			)},
			"clientid":   {this.ClientId},
			"psessionid": {"null"},
		})
	if err != nil {
		return
	}
	hr = &Login2Result{}
	err = json.NewDecoder(res.Body).Decode(hr)
	res.Body.Close()
	if err != nil {
		return
	}
	return
}

// 登录
func (this *WebQQ) Login() (err error) {
	defer util.Catch()
	this.LoginSig, err = this.ptlogin_login_sig()
	DEBUG.Log("[login_sig]", this.LoginSig)
	util.Try(err)
	_, code, pwd, err := this.ptlogin_check()
	util.Try(err)
	DEBUG.Log("[ptlogin_check] RET OK ", code, pwd)
	err = this.ptlogin_login(code, pwd)
	util.Try(err)
	DEBUG.Log("[ptlogin_login] RET OK ")
	ptwebqq := this.GetCookie(util.MustParseUrl(PTLOGIN_URL), "ptwebqq")
	if ptwebqq == "" {
		return fmt.Errorf("[ptwebqq] Failed to read cookie.")
	}
	ret, err := this.channel_login2(ptwebqq)
	util.Try(err)
	if ret.Code != 0 {
		if ret.Code == 103 {
			ret.Msg = "Error 103"
		}
		return fmt.Errorf("%v : %s\n%v", ret.Code, ret.Msg, ptwebqq)
	}
	this.VerifyCode = ret.Result.VerifyCode
	this.SessionId = ret.Result.SessionId
	INFO.Log("Login success")
	return
}

func (this *WebQQ) GetCookie(url *url.URL, name string) (ret string) {
	if this.C.Jar != nil {
		for _, v := range this.C.Jar.Cookies(url) {
			DEBUG.Log(v)
			if v.Name == name {
				ret = v.Value
				// return
			}
		}
	}
	return
}

// 带变量Referer GET
func (this *WebQQ) GetWithReferer(urlStr, referer string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if this.C.Jar != nil {
		for _, v := range this.C.Jar.Cookies(req.URL) {
			req.AddCookie(v)
		}
	}
	req.Header.Add("Referer", referer)
	return this.C.Do(req)
}

// 带固定Referer POST
func (this *WebQQ) PostFormWithReferer(url string, val url.Values) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(val.Encode()))
	if this.C.Jar != nil {
		for _, v := range this.C.Jar.Cookies(req.URL) {
			req.AddCookie(v)
		}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://s.web2.qq.com/proxy.html?v=20110412001&callback=1&id=1")
	return this.C.Do(req)
}

// 进行三次加盐MD5之类...什么的算法
func (this *WebQQ) GenPwd(vSalt, vCode string) string {
	vSaltedPwd := util.Md5String(this.PasswdMd5 + vSalt)
	return util.Md5StringX(strings.ToUpper(vSaltedPwd + vCode))
}

// ptlogin_login的返回值 JSON
type Login2Result struct {
	Code   float64 `json:"retcode"`
	Msg    string  `json:"errmsg"`
	Result struct {
		ClientIP   uint32 `json:"cip"`
		F          uint32 `json:"f"`
		Index      uint32 `json:"index"`
		Port       uint32 `json:"port"`
		SessionId  string `json:"psessionid"`
		Status     string `json:"status"`
		Uin        uint64 `json:"uin"`
		UserState  uint32 `json:"user_state"`
		VerifyCode string `json:"vfwebqq"`
	} `json:"result"`
}
