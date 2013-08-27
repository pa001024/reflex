package webqq

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/pa001024/MoeWorker/util"
)

// WebQQ对象
type WebQQ struct {
	client *http.Client

	Id        Account
	PasswdMd5 string
	// 瞬态
	IdStr      string
	Uin        Uin
	ClientId   string
	VerifyCode string
	SessionId  string
	PtWebQQ    string
	LoginSig   string
	SigUrl     string
}

// 用户ID
type Uin uint64

func (u Uin) String() string { return fmt.Sprint(uint64(u)) }

// 群信息ID
type GCode uint64

func (u GCode) String() string { return fmt.Sprint(uint64(u)) }

// 真实QQ号
type Account uint64

func (u Account) String() string { return fmt.Sprint(uint64(u)) }

// 创建WebQQ
func NewWebQQ(uid Account, pwd string) (this *WebQQ) {
	jar, _ := cookiejar.New(nil)
	this = &WebQQ{
		client:    &http.Client{nil, nil, jar},
		Id:        uid,
		IdStr:     util.ToString(uid),
		PasswdMd5: pwd,
		ClientId:  fmt.Sprint(rand.Int31n(90000000) + 10000000),
		SigUrl:    fmt.Sprintf("https://ui.ptlogin2.qq.com/cgi-bin/login?daid=164&target=self&style=5&mibao_css=m_webqq&appid=1003903&enable_qlogin=0&no_verifyimg=1&s_url=http%%3A%%2F%%2Fweb2.qq.com%%2Floginproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=%s001", time.Now().Format("20060102")),
	}
	return
}

// 登录 [2013.8.27]
func (this *WebQQ) Login() (err error) {
	defer util.Catch()
	// [1]
	this.LoginSig, err = this.ptlogin_login_sig()
	util.DEBUG.Logf("[ptlogin_login_sig] login_sig = %s ", this.LoginSig)
	// [2]
	_, code, pwd, err := this.ptlogin_check()
	util.DEBUG.Logf("[ptlogin_check] Return OK %s %s", code, pwd)
	// [3]
	pturl, msg, err := this.ptlogin_login(code, pwd)
	util.DEBUG.Logf("[ptlogin_login] Return %s and check_sig = %s", msg, pturl)
	// [4]
	err = this.ptlogin_check_sig(pturl)
	if this.PtWebQQ = this.getCookie(util.MustParseUrl(PTLOGIN_URL), "ptwebqq"); this.PtWebQQ == "" {
		return fmt.Errorf("[ptwebqq] Failed to read cookie.")
	}
	// [5]
	ret, err := this.login2()
	if ret.Code != 0 {
		return fmt.Errorf("[channel_login2] %v : %s", ret.Code, ret.Msg)
	}
	this.VerifyCode = ret.Result.VerifyCode
	this.SessionId = ret.Result.SessionId
	this.Uin = ret.Result.Uin
	util.INFO.Log("[Login] Login success")
	return
}

var (
	msg_id uint32 = (2000 + uint32(rand.Int31n(2999))) * 1000
)

// 给Uin发送消息
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
				util.WARN.Logf("poll2() throw error: %v", err)
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
