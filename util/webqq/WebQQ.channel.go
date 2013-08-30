package webqq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/pa001024/MoeWorker/util"
)

const (
	_CHANNEL_URL     = "http://d.web2.qq.com/channel/"
	_CHANNEL_REFERER = "http://d.web2.qq.com/proxy.html?v=20110331002&callback=1&id=2"
)

// 通用channel接口(GET)
func (this *WebQQ) channel(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"clientid":   {this.clientid},
		"psessionid": {this.psessionid},
		"t":          {util.JsCurrentTime()},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.getWithReferer(fmt.Sprintf("%s%s?%s", _CHANNEL_URL, api, val.Encode()), _CHANNEL_REFERER)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用channel接口(POST)
func (this *WebQQ) postChannel(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"r": {util.ToJson(
			append(args, "clientid", this.clientid, "psessionid", this.psessionid)...,
		)},
		"clientid":   {this.clientid},
		"psessionid": {this.psessionid},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.postFormWithReferer(_CHANNEL_URL+api, _CHANNEL_REFERER, val)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 登录
func (this *WebQQ) login2() (v *ResultLogin2, err error) {
	data, err := this.postChannel("login2", "status", "online", "ptwebqq", this.ptwebQQ, "passwd_sig", "")
	if err == nil {
		v = &ResultLogin2{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 登录 返回值结构
type ResultLogin2 struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result struct {
		Uin        Uin    `json:"uin"`
		VerifyCode string `json:"vfwebqq"`
		SessionId  string `json:"psessionid"`
		Status     string `json:"status"`
	} `json:"result"`
}

// 获取消息
func (this *WebQQ) poll2() (v *ResultPoll, err error) {
	data, err := this.postChannel("poll2")
	if err == nil {
		v = &ResultPoll{}
		err = json.Unmarshal(data, v)
	}
	return
}

// poll2 result结构
type ResultPoll struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result []struct {
		Type  string          `json:"poll_type"`
		Value json.RawMessage `json:"value"` // see WebQQ.event.go
	} `json:"result"`
}

// 获取在线好友
func (this *WebQQ) get_online_buddies2() (v *ResultOnlineBuddies, err error) {
	data, err := this.channel("get_online_buddies2")
	if err == nil {
		v = &ResultOnlineBuddies{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 在线好友result结构
type ResultOnlineBuddies struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin        Uin    `json:"uin"`
		Status     string `json:"status"`
		ClientType uint32 `json:"client_type"`
	} `json:"result"`
}

// 正在输入
func (this *WebQQ) input_notify2(to_uin string) (v *Result, err error) {
	data, err := this.channel("input_notify2", "to_uin", to_uin)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 窗口抖动
func (this *WebQQ) shake2(to_uin string) (v *Result, err error) {
	data, err := this.channel("shake2", "to_uin", to_uin)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取当前聊天列表
func (this *WebQQ) get_recent_list2() (v *ResultRecentList, err error) {
	data, err := this.postChannel("get_recent_list2")
	if err == nil {
		v = &ResultRecentList{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 当前聊天列表result结构 {"retcode":0,"result":[{"uin":3255435951,"type":0},{"uin":221664830,"type":1}]}
type ResultRecentList struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin  Uin    `json:"uin"`
		Type uint32 `json:"type"`
	} `json:"result"`
}

// 发送私聊消息
func (this *WebQQ) send_buddy_msg2(to Uin, content ContentModel, msg_id uint32) (v *Result, err error) {
	util.DEBUG.Logf("send_buddy_msg2(to = %s , content = %v , msg_id = %v)", to, content, msg_id)
	data, err := this.postChannel("send_buddy_msg2",
		"to", to,
		"face", "552", // 这是啥?
		"content", content.Encode(this).EncodeString(),
		"msg_id", msg_id,
	)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 发送群消息
func (this *WebQQ) send_qun_msg2(group_uin Uin, content ContentModel, msg_id uint32) (v *Result, err error) {
	util.DEBUG.Logf("send_qun_msg2(group_uin = %s , content = %v , msg_id = %v)", group_uin, content, msg_id)
	data, err := this.postChannel("send_qun_msg2",
		"group_uin", group_uin,
		"content", content.Encode(this).EncodeString(),
		"msg_id", msg_id,
	)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取发送群聊图片Sig
func (this *WebQQ) get_gface_sig2() (v *ResultGfaceSig, err error) {
	data, err := this.postChannel("get_gface_sig2")
	if err == nil {
		v = &ResultGfaceSig{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 发送群聊图片Sig result结构
type ResultGfaceSig struct {
	Code   int `json:"retcode"`
	Result struct {
		Reply    uint32 `json:"reply"`
		GfaceKey string `json:"gface_key"`
		GfaceSig string `json:"gface_sig"`
	} `json:"result"`
}

// 发送离线图片
func (this *WebQQ) apply_offline_pic_dl2(file_path string) (v *ResultOfflinePicDl, err error) {
	data, err := this.channel("apply_offline_pic_dl2",
		"f_uin", this.Uin,
		"file_path", file_path,
	)
	if err == nil {
		v = &ResultOfflinePicDl{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 发送离线图片 result结构
type ResultOfflinePicDl struct {
	Code   int `json:"retcode"`
	Result struct {
		Success  uint32 `json:"success"`
		url      string `json:"url"`
		FilePath string `json:"file_path"`
	} `json:"result"`
}

// 创建讨论组
func (this *WebQQ) create_discu(file_path string) (v *ResultCreateDiscu, err error) {
	data, err := this.channel("create_discu") // TODO

	if err == nil {
		v = &ResultCreateDiscu{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 创建讨论组 result结构
type ResultCreateDiscu struct {
	Code   int `json:"retcode"`
	Result struct {
		Result  uint32  `json:"result"`
		DiscuId DiscuId `json:"did"`
	} `json:"result"`
}

// 获取讨论组信息
func (this *WebQQ) get_discu_info(did DiscuId) (v *ResultDiscuInfo, err error) {
	data, err := this.channel("get_discu_info",
		"did", did,
		"vfwebqq", this.vfwebqq,
	)
	if err == nil {
		v = &ResultDiscuInfo{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取讨论组信息 result结构
type ResultDiscuInfo struct {
	Code   int `json:"retcode"`
	Result struct {
		Info struct {
			DiscuId    DiscuId `json:"did"`
			DiscuOwner Uin     `json:"discu_owner"`
			DiscuName  string  `json:"discu_name"`
			InfoSeq    string  `json:"info_seq"` // 信息序列号
			MemberList []struct {
				Uin     Uin     `json:"mem_uin"`
				Account Account `json:"ruin"`
			} `json:"mem_list"`
		} `json:"info"`
		MemberStatus []MemberStat `json:"mem_status"`
		MemberInfo   struct {
			Uin  Uin    `json:"uin"`
			Nick string `json:"nick"`
		} `json:"mem_info"`
	} `json:"result"`
}

// 获取讨论组信息
func (this *WebQQ) modify_discu_info(did DiscuId, discu_name string) (v *Result, err error) {
	data, err := this.postChannel("modify_discu_info",
		"did", did,
		"discu_name", discu_name,
		"dtype", 1,
		"vfwebqq", this.vfwebqq,
	)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取临时会话Sig
func (this *WebQQ) get_c2cmsg_sig2(to_uin Uin) (v *ResultC2CMsgSig, err error) {
	data, err := this.channel("get_c2cmsg_sig2",
		"id", this.Uin,
		"to_uin", to_uin,
		"service_type", 1,
	)
	if err == nil {
		v = &ResultC2CMsgSig{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取临时会话Sig result结构
type ResultC2CMsgSig struct {
	Code   int `json:"retcode"`
	Result struct {
		Type  uint32 `json:"type"`
		Value string `json:"value"`
		Flags struct {
			Text  uint8 `json:"text"`
			Pic   uint8 `json:"pic"`
			File  uint8 `json:"file"`
			Audio uint8 `json:"audio"`
			Video uint8 `json:"video"`
		} `json:"flags"`
	} `json:"result"`
}
