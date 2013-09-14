package webqq

import (
	"encoding/json"
	"fmt"

	"github.com/pa001024/reflex/util"
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
	case "input_notify":
		v = &EventInputNotify{}
	case "filesrv_transfer":
		v = &EventFileSrvTransfer{}
	case "file_message":
		v = &EventFileMessage{}
	case "av_refuse":
		v = &EventAvRefuse{}
	case "shake_message":
		v = &EventShakeMessage{}
	case "push_offfile":
		v = &EventPushOffFile{}
	case "system_message":
		v = &EventSystemMessage{}
	case "sys_g_msg":
		v = &EventSysGroupMessage{}
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

// "poll_type": "av_refuse"
type EventAvRefuse struct {
	FromUin   Uin    `json:"from_uin"`
	ToUin     Uin    `json:"to_uin"`
	MsgId     uint32 `json:"msg_id"`
	MsgId2    uint32 `json:"msg_id2"`
	MsgType   uint32 `json:"msg_type"`
	ReplyIp   uint32 `json:"reply_ip"`
	Type      uint32 `json:"type"`
	SessionId uint32 `json:"session_id"`
	Time      uint32 `json:"time"`
	Longconn  string `json:"longconn"`
}

// "poll_type": "shake_message"
type EventShakeMessage struct {
	FromUin Uin    `json:"from_uin"`
	ToUin   Uin    `json:"to_uin"`
	MsgId   uint32 `json:"msg_id"`
	MsgId2  uint32 `json:"msg_id2"`
	MsgType uint32 `json:"msg_type"`
	ReplyIp uint32 `json:"reply_ip"`
}

// "poll_type": "push_offfile"
type EventPushOffFile struct {
	FromUin    Uin    `json:"from_uin"`
	MsgId      uint32 `json:"msg_id"`
	Rkey       string `json:"rkey"`
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	Port       uint32 `json:"port"`
	Size       uint32 `json:"size"`
	ExpireTime uint32 `json:"expire_time"`
	Time       uint32 `json:"time"`
}

type SystemMessageType string

const (
	AddedBuddySig   SystemMessageType = "added_buddy_sig"
	AddedBuddyNosig                   = "added_buddy_nosig"
	VerifyPassAdd                     = "verify_pass_add"
	VerifyPass                        = "verify_pass"
	VerifyRequired                    = "verify_required"
	VerifyRejected                    = "verify_rejected"
)

// "poll_type": "system_message"
type EventSystemMessage struct {
	Type       SystemMessageType `json:"type"`
	FromUin    Uin               `json:"from_uin"`
	Stat       uint32            `json:"stat"`
	ClientType ClientType        `json:"client_type"`
	GroupId    GCode             `json:"group_id"`
	Msg        string            `json:"msg"`
	Account    Account           `json:"account"`
	Uiuin      Uin               `json:"uiuin"`
}

type SysGroupMessageType string

const (
	GroupJoin             SysGroupMessageType = "group_join"
	GroupLeave                                = "group_leave"
	GroupRequestJoin                          = "group_request_join"
	GroupRequestJoinAgree                     = "group_request_join_agree"
	GroupRequestJoinDeny                      = "group_request_join_deny"
	GroupAdminOp                              = "group_admin_op"
	GroupCreate                               = "group_create"
)

// "poll_type": "sys_g_msg"
type EventSysGroupMessage struct {
	Type       SysGroupMessageType `json:"type"`
	NewMember  string              `json:"new_member"`
	OldMember  string              `json:"old_member"`
	GCode      GCode               `json:"gcode"`
	RequestUin Uin                 `json:"request_uin"`
	OpType     string              `json:"op_type"` // 0,2:取消管理员 3:成为管理员 255:群主(转让)
	// TOldMember  string `json:"t_old_member"`
	// TGCode      GCode  `json:"t_gcode"`
	// TNewMember  string `json:"t_new_member"`
	// TRequestUin Uin    `json:"t_request_uin"`
}
