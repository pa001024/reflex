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
	Uin        Uin
	ClientId   string
	VerifyCode string
	SessionId  string
	PtWebQQ    string
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
	this.PtWebQQ = this.GetCookie(util.MustParseUrl(PTLOGIN_URL), "ptwebqq")
	if this.PtWebQQ == "" {
		return fmt.Errorf("[ptwebqq] Failed to read cookie.")
	}
	ret, err := this.channel_login2()
	util.Try(err)
	if ret.Code != 0 {
		if ret.Code == 103 {
			ret.Msg = "Error 103"
		}
		return fmt.Errorf("%v : %s\n%v", ret.Code, ret.Msg, this.PtWebQQ)
	}
	this.VerifyCode = ret.Result.VerifyCode
	this.SessionId = ret.Result.SessionId
	this.Uin = ret.Result.Uin
	INFO.Log("Login success")
	return
}

// 进行三次加盐MD5之类...什么的算法 [2013.8.6]
func (this *WebQQ) GenPwd(salt, code string) string {
	vSaltedPwd := util.Md5String(this.PasswdMd5 + salt)
	return util.Md5StringX(vSaltedPwd + strings.ToUpper(code))
}

var (
	msg_id uint32 = (2000 + uint32(rand.Int31n(2999))) * 1000
)

func (this *WebQQ) SendTo(to Uin, m ContentM) (err error) {
	r, err := this.send_buddy_msg2(to, m, msg_id)
	msg_id++
	if r != nil && err == nil {
		err = fmt.Errorf("SendTo() return code %v", r.Code)
	}
	return
}

// 开始接受消息并发送到channel
func (this *WebQQ) Start() <-chan Event {
	in := make(chan *ResultPoll, 3) // 防止被消息处理阻塞 可调大
	out := make(chan Event)
	go func() {
		for {
			r, err := this.poll2() // TODO: 超时处理
			if err != nil {
				WARN.Logf("poll2() throw error: %v", err)
			} else if r != nil {
				in <- r
			}
		}
	}()
	for {
		for _, v := range (<-in).Result {
			e, err := RawEvent(v.Value).ParseEvent(v.Type)
			if err == nil {
				out <- e
			}
		}
	}
	return out
}
