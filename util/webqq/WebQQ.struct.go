package webqq

// ptlogin_login的返回值 JSON
// 样本:
/*
{
	"retcode": 0,
	"result": {
		"uin": 2735284921,
		"cip": 3080236829,
		"index": 1075,
		"port": 47943,
		"status": "online",
		"vfwebqq": "209da35a9665546efac6e1032577fd75e8fcae3e2d7a264fc64fe598064245285ae63a270bc204f4",
		"psessionid": "8368046764001d636f6e6e7365727665725f77656271714031302e3133392e372e3136300000443100000163026e0400b92209a36d0000000a404b454773376a7457646d00000028209da35a9665546efac6e1032577fd75e8fcae3e2d7a264fc64fe598064245285ae63a270bc204f4",
		"user_state": 0,
		"f": 0
	}
}
*/
type Login2Result struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result struct {
		Uin        uint64 `json:"uin"`
		VerifyCode string `json:"vfwebqq"`
		SessionId  string `json:"psessionid"`
		Status     string `json:"status"`
		// CIP        uint32 `json:"cip"` // 没用
		// Index     uint32 `json:"index"` // 没用
		// Port      uint32 `json:"port"`       // 没用
		// UserState uint32 `json:"user_state"` // 没用
		// F          uint32 `json:"f"` // 没用
	} `json:"result"`
}

// poll2 返回值
// 样本:
/*
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
				["font",
				{
					"size": 10,
					"color": "000000",
					"style": [0, 0, 0],
					"name": "\u5FAE\u8F6F\u96C5\u9ED1"
				}], "123123 "]
		}
	}]
}
*/
type Poll2Result struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result []struct {
		Type  string  `json:"poll_type"`
		Value Message `json:"value"`
	} `json:"result"`
}

// 消息结构
type Message struct {
	Id      uint32 `json:"msg_id"`   // 消息Id 防止重复回复
	Id2     uint32 `json:"msg_id2"`  // 同上
	Type    uint32 `json:"msg_type"` // 消息类型
	ReplyIP uint32 `json:"reply_ip"` // 返回号
	From    uint64 `json:"from_uin"` // 独立群号/用户号
	// To      uint64 `json:"to_uin"`// 没用
	Time uint64 `json:"time"` // 发送时间
	// GroupCode uint64 `json:"group_code"`// 没用
	Sender uint64 `json:"send_uin"` // 独立用户号
	// Seq       uint64   `json:"seq"` // 没用
	GroupId uint64   `json:"info_seq"` // 明文群号
	Content []string `json:"content"`  // Content[0] 本来是字体 现在直接忽略掉
}

func (this *Message) Filter() {
	arr := make([]string, 0, len(this.Content)-1)
	for _, v := range this.Content {
		if v != "" {
			arr = append(arr, v)
		}
	}
	this.Content = arr
}

// 发送私聊消息结构 send_buddy_msg2
// 样本:
/*
{
	"to": 3255435951,
	"face": 552,
	"content": "[\"asd\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
	"msg_id": 38610005,
	"clientid": "10861648",
	"psessionid": "8368046764001d636f6e6e7365727665725f77656271714031302e3133392e372e3136300000443100000163026e0400b92209a36d0000000a404b454773376a7457646d00000028209da35a9665546efac6e1032577fd75e8fcae3e2d7a264fc64fe598064245285ae63a270bc204f4"
}
*/
type SendMessage struct {
	To        uint64 `json:"to"`
	Face      uint32 `json:"face"`
	Content   string `json:"content"`
	MessageId uint32 `json:"msg_id"`
	ClientId  string `json:"clientid"`
	SessionId string `json:"psessionid"`
}

// 发送群消息结构 send_qun_msg2
// 样本:
/*
{
	"group_uin": 221664830,
	"content": "[\"123123123\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
	"msg_id": 38610004,
	"clientid": "10861648",
	"psessionid": "8368046764001d636f6e6e7365727665725f77656271714031302e3133392e372e3136300000443100000163026e0400b92209a36d0000000a404b454773376a7457646d00000028209da35a9665546efac6e1032577fd75e8fcae3e2d7a264fc64fe598064245285ae63a270bc204f4"
}
*/
type SendQunMessage struct {
	To        uint64 `json:"group_uin"`
	Content   string `json:"content"`
	MessageId uint32 `json:"msg_id"`
	ClientId  string `json:"clientid"`
	SessionId string `json:"psessionid"`
}
