package webqq

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/pa001024/MoeWorker/util"
)

const CHANNEL_URL = "http://d.web2.qq.com/channel/"

// 通用channel接口(GET)
func (this *WebQQ) channel(api string, args ...string) (body []byte, err error) {
	val := url.Values{
		"clientid":   {this.ClientId},
		"psessionid": {this.SessionId},
		"t":          {fmt.Sprintf("%d", time.Now().UnixNano()/1e6)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i], args[i+1])
	}
	res, err := this.GetWithReferer(fmt.Sprintf("%s%s?%s", CAPTCHA_URL, api, val.Encode()))
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用channel接口(POST)
func (this *WebQQ) postChannel(api string, args ...string) (body []byte, err error) {
	val := url.Values{
		"clientid":   {this.ClientId},
		"psessionid": {this.SessionId},
		"t":          {fmt.Sprintf("%d", time.Now().UnixNano()/1e6)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i], args[i+1])
	}
	res, err := this.PostFormWithReferer(CAPTCHA_URL+api, val)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

/*
获取在线好友
------------

{"retcode":0,"result":[{"uin":3255435951,"status":"online","client_type":1}]}
*/
func (this *WebQQ) get_online_buddies2() {
	this.channel("get_online_buddies2")
}

/*
正在输入
--------

{"retcode":0,"result":"ok"}
*/
func (this *WebQQ) input_notify2(to_uin string) {
	this.channel("input_notify2", "to_uin", to_uin)
}

/*
窗口抖动
--------

{"retcode":0,"result":"ok"}
*/
func (this *WebQQ) shake2(to_uin string) {
	this.channel("shake2", "to_uin", to_uin)
}

/*
获取
----

{"retcode":0,"result":[{"uin":3255435951,"type":0},{"uin":221664830,"type":1}]}
*/
func (this *WebQQ) get_recent_list2() {
	this.postChannel("get_recent_list2")
}

/*
发送私聊消息
------------

r = {
	"to": 3255435951,
	"face": 552,
	"content": "[\"asd\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
	"msg_id": 38610005,
	"clientid": "10861648",
	"psessionid": "..."
}
{"retcode":0,"result":"ok"}
*/
func (this *WebQQ) send_buddy_msg2() {
	this.postChannel("send_buddy_msg2")
}

// send_buddy_msg2 r结构
type SendBuddyMessage struct {
	To        Uin    `json:"to"`
	Face      uint32 `json:"face"` // 这是啥?
	Content   string `json:"content"`
	MessageId uint32 `json:"msg_id"`
	ClientId  string `json:"clientid"`
	SessionId string `json:"psessionid"`
}

/*
发送群消息
----------

r = {
	"group_uin": 221664830,
	"content": "[\"msg\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
	"msg_id": 38610004,
	"clientid": "10861648",
	"psessionid": "..."
}
{"retcode":0,"result":"ok"}
*/
func (this *WebQQ) send_qun_msg2(msg string) {
	this.postChannel("send_qun_msg2", r, util.ToJson(
		"r", "",
	))
}

// send_qun_msg2 r结构
type SendQunMessage struct {
	To        Uin    `json:"group_uin"`
	Content   string `json:"content"`
	MessageId uint32 `json:"msg_id"`
	ClientId  string `json:"clientid"`
	SessionId string `json:"psessionid"`
}

// 获取消息
func (this *WebQQ) poll2() {
	this.postChannel("poll2")
}

/*
poll2 返回值
------------

{"retcode":116,"p":"39bd5c71be123aaf451073154d52bfef78b1adaa0d087601"} // 什么都没有
{
	"retcode": 0,
	"result": [{
		"poll_type": "message",
		"value": {
			"msg_id": 31607,
			"from_uin": 3255435951,
			"to_uin": 2735284921,
			"msg_id2": 459162,
			"msg_type": 9,
			"reply_ip": 176756886,
			"time": 1375794494,
			"content": [
					["font",
					{
						"size": 10,
						"color": "000000",
						"style": [0, 0, 0],
						"name": "\u5FAE\u8F6F\u96C5\u9ED1"
					}],
					"\u53E5\u53E5\u53E5 "
				]
		}
	},{
		"poll_type": "group_message",
		"value": {
			"msg_id": 12418,
			"from_uin": 221664830,
			"to_uin": 2735284921,
			"msg_id2": 7256,
			"msg_type": 43,
			"reply_ip": 176496859,
			"group_code": 738328699,
			"send_uin": 3255435951,
			"seq": 114,
			"time": 1375798045,
			"info_seq": 165640562,
			"content": [
				["font",{
					"size": 10,
					"color": "000000",
					"style": [0, 0, 0],
					"name": "\u5FAE\u8F6F\u96C5\u9ED1"
				}],
				"123123 "
			]
		}
	},{
		"poll_type": "buddies_status_change",
		"value": {
			"uin": 3255435951,
			"status": "online",
			"client_type": 1
		}
	}]
}
*/
type ResultPoll struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result []struct {
		Type  string `json:"poll_type"`
		Value Poll   `json:"value"`
	} `json:"result"`
}

// poll2 消息结构
type Poll struct {
	Id         uint32   `json:"msg_id"`      // 消息Id 防止重复回复
	Id2        uint32   `json:"msg_id2"`     // 同上
	Type       uint32   `json:"msg_type"`    // 消息类型
	ReplyIP    uint32   `json:"reply_ip"`    // 返回号
	From       uint64   `json:"from_uin"`    // 独立群号/用户号
	To         uint64   `json:"to_uin"`      // 自己的UIN
	Time       uint64   `json:"time"`        // 发送时间
	GroupCode  uint64   `json:"group_code"`  // gcode
	Sender     uint64   `json:"send_uin"`    // 独立用户号
	Seq        uint64   `json:"seq"`         // 没用
	GroupId    uint64   `json:"info_seq"`    // 明文群号
	Status     string   `json:"status"`      // 状态变化
	ClientType uint32   `json:"client_type"` // 客户端类型
	Content    []string `json:"content"`     // Content[0] 本来是字体 现在直接忽略掉
}

func (this *Poll) Filter() {
	arr := make([]string, 0, len(this.Content)-1)
	for _, v := range this.Content {
		if v != "" {
			arr = append(arr, v)
		}
	}
	this.Content = arr
}
