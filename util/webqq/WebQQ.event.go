package webqq

import (
	"encoding/json"
	"fmt"

	"github.com/pa001024/MoeWorker/util"
)

// 解析事件
func (data RawEvent) ParseEvent(poll_type string) (v Event, err error) {
	util.TRACE.Logf("[RawEvent] %s", string(data))
	switch poll_type {
	case "message":
		v = &EventMessage{}
	case "group_message":
		v = &EventGroupMessage{}
	case "discu_message":
		v = &EventDiscuMessage{}
	case "buddies_status_change":
		v = &EventBuddiesStatusChange{}
	case "group_web_message":
		v = &EventGroupWebMessage{}
	default:
		err = fmt.Errorf("Unsupport poll_type: %v", poll_type)
		return
	}
	err = json.Unmarshal(data, v)
	return
}

type RawEvent json.RawMessage
type Event interface{}

// "poll_type": "message"
type EventMessage struct {
	FromUin Uin     `json:"from_uin"`
	ToUin   Uin     `json:"to_uin"`
	MsgId   uint32  `json:"msg_id"`
	MsgId2  uint32  `json:"msg_id2"`
	MsgType uint32  `json:"msg_type"`
	ReplyIp uint32  `json:"reply_ip"`
	Time    uint64  `json:"time"`
	Content Content `json:"content"`
}

// "poll_type": "group_message"
type EventGroupMessage struct {
	FromUin   Uin     `json:"from_uin"`
	ToUin     Uin     `json:"to_uin"`
	MsgId     uint32  `json:"msg_id"`
	MsgId2    uint32  `json:"msg_id2"`
	MsgType   uint32  `json:"msg_type"`
	ReplyIp   uint32  `json:"reply_ip"`
	GroupCode GCode   `json:"group_code"`
	SendUin   Uin     `json:"send_uin"`
	Seq       uint64  `json:"seq"`
	Time      uint64  `json:"time"`
	InfoSeq   uint64  `json:"info_seq"` // 群号
	Content   Content `json:"content"`
}

// "poll_type": "discu_message"
type EventDiscuMessage struct {
	FromUin Uin     `json:"from_uin"`
	ToUin   Uin     `json:"to_uin"`
	MsgId   uint32  `json:"msg_id"`
	MsgId2  uint32  `json:"msg_id2"`
	MsgType uint32  `json:"msg_type"`
	ReplyIp uint32  `json:"reply_ip"`
	DiscuId DiscuId `json:"did"`
	SendUin Uin     `json:"send_uin"`
	Seq     uint64  `json:"seq"`
	Time    uint64  `json:"time"`
	InfoSeq uint64  `json:"info_seq"` // 讨论组号
	Content Content `json:"content"`
}

// "poll_type": "buddies_status_change"
type EventBuddiesStatusChange struct {
	Uin        Uin    `json:"uin"`
	Status     uint32 `json:"status"`
	ClientType uint32 `json:"client_type"`
}

// "poll_type": "group_web_message"
type EventGroupWebMessage struct {
	FromUin   Uin    `json:"from_uin"`
	ToUin     Uin    `json:"to_uin"`
	MsgId     uint32 `json:"msg_id"`
	MsgId2    uint32 `json:"msg_id2"`
	MsgType   uint32 `json:"msg_type"`
	ReplyIp   uint32 `json:"reply_ip"`
	GroupType uint32 `json:"group_type"`
	Ver       uint32 `json:"ver"`
	GroupCode GCode  `json:"group_code"`
	SendUin   Uin    `json:"send_uin"`
	Xml       string `json:"xml"`
}

// "poll_type": "input_notify"
type EventInputNotify struct {
	FromUin Uin    `json:"from_uin"`
	ToUin   Uin    `json:"to_uin"`
	MsgId   uint32 `json:"msg_id"`
	MsgId2  uint32 `json:"msg_id2"`
	MsgType uint32 `json:"msg_type"`
	ReplyIp uint32 `json:"reply_ip"`
}

// "poll_type": "filesrv_transfer"
type EventFileSrvTransfer struct {
	FromUin   Uin    `json:"from_uin"`
	ToUin     Uin    `json:"to_uin"`
	LcId      uint32 `json:"lc_id"`
	Type      uint32 `json:"type"`
	Operation uint32 `json:"operation"`
	Now       uint32 `json:"now"`
	FileCount uint32 `json:"file_count"`
	FileInfos struct {
		FileName   string `json:"file_name"`
		ProId      uint32 `json:"pro_id"`
		FileStatus uint32 `json:"file_status"`
	} `json:"file_infos"`
}

// "poll_type": "file_message"
type EventFileMessage struct {
	FromUin   Uin    `json:"from_uin"`
	ToUin     Uin    `json:"to_uin"`
	Mode      string `json:"mode"`
	MsgId     uint32 `json:"msg_id"`
	MsgId2    uint32 `json:"msg_id2"`
	MsgType   uint32 `json:"msg_type"`
	ReplyIp   uint32 `json:"reply_ip"`
	Type      uint32 `json:"type"`
	SessionId uint32 `json:"session_id"`
	Time      uint32 `json:"time"`
	InetIp    uint32 `json:"inet_ip"`
}
