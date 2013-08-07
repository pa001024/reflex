package webqq

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
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

// 通用channel接口
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

// 获取在线好友
/*
{"retcode":0,"result":[]}
*/
func (this *WebQQ) get_online_buddies2() {
	this.channel("get_online_buddies2")
}

// 获取
/*
{"retcode":0,"result":[{"uin":3255435951,"type":0},{"uin":221664830,"type":1}]}
*/
func (this *WebQQ) get_recent_list2() {
	this.postChannel("get_recent_list2")
}

// 获取消息
/*
{"retcode":116,"p":"39bd5c71be123aaf451073154d52bfef78b1adaa0d087601"} // 什么都没有
*/
func (this *WebQQ) poll2() {
	this.postChannel("poll2")
}
