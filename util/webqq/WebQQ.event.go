package webqq

import (
	"encoding/json"
	"fmt"

	"github.com/pa001024/MoeWorker/util"
)

// 解析事件
func (data RawEvent) ParseEvent(poll_type string) (v Event, err error) {
	switch poll_type {
	case "message":
		d := &EventMessage{}
		err = json.Unmarshal(data, d)
		v = d
	case "group_message":
		d := &EventGroupMessage{}
		err = json.Unmarshal(data, d)
		v = d
	case "buddies_status_change":
		d := &EventBuddiesStatusChange{}
		err = json.Unmarshal(data, d)
		v = d
	default:
		err = fmt.Errorf("Unsupport poll_type: %v", poll_type)
		util.DEBUG.Log(string(data))
	}
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

// "poll_type": "buddies_status_change"
type EventBuddiesStatusChange struct {
	Uin        Uin    `json:"uin"`
	Status     uint32 `json:"status"`
	ClientType uint32 `json:"client_type"`
}
