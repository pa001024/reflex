package webqq

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/pa001024/MoeWorker/util"
)

// 三个日志级别
var (
	DEBUG = util.NewLogger(false, "[DEBUG] ")
	WARN  = util.NewLogger(true, "[WARNING] ")
	INFO  = util.NewLogger(true, "[INFO] ")
)

// WebQQ对象
type WebQQ struct {
	client *http.Client

	Id        string
	PasswdMd5 string
	// 瞬态
	ClientId   string
	VerifyCode string
	SessionId  string
}

// 创建WebQQ
func NewWebQQ(uid, pwd string) (this *WebQQ) {
	jar, _ := cookiejar.New(nil)
	this = &WebQQ{
		client:    &http.Client{nil, nil, jar},
		Id:        uid, // 可以是邮箱或者QQ号
		PasswdMd5: pwd,
		ClientId:  fmt.Sprint(rand.Int31n(90000000) + 10000000),
	}
	return
}

// 登录 [2013.8.6]
func (this *WebQQ) Login() (err error) {
	defer util.Catch()
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

// 进行三次加盐MD5之类...什么的算法 [2013.8.6]
func (this *WebQQ) GenPwd(salt, code string) string {
	vSaltedPwd := util.Md5String(this.PasswdMd5 + salt)
	return util.Md5StringX(vSaltedPwd + strings.ToUpper(code))
}

//
func (this *WebQQ) SendTo(qid string) {

}

// 开始接受消息并发送到channel
func (this *WebQQ) Start() {
	// make(chan )
}
