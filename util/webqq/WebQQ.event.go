package webqq

import (
	"encoding/json"
	"fmt"
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
	}
	return
}

type RawEvent json.RawMessage
type Event interface{}

/*
 "poll_type": "message"
 {
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
*/
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

/*
 "poll_type": "group_message"
 {
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
*/
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

/*
 "poll_type": "buddies_status_change"
 {
 	"uin": 3255435951,
 	"status": "online",
 	"client_type": 1
 }
*/
type EventBuddiesStatusChange struct {
	Uin        Uin    `json:"uin"`
	Status     uint32 `json:"status"`
	ClientType uint32 `json:"client_type"`
}
