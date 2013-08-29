/*
 登录 返回值结构
 ---------------
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
 在线好友result结构
 ------------------

 {"retcode":0,"result":[{"uin":3255435951,"status":"online","client_type":1}]}

 发送私聊消息
 ------------

 to:3255435951,
 face:552,
 content:"[\"asd\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
 msg_id:38610005,
 clientid:"10861648",
 psessionid:"..."
 {"retcode":0,"result":"ok"}

 发送群消息
 ----------

 r = {
 	"group_uin": 221664830,
 	"content": "[\"msg\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
 	"msg_id": 38610004,
 	"clientid": "10861648",
 	"psessionid": "..."
 }

 poll2 result结构
 ----------------
 {"retcode":102,"errmsg":""} // 什么都没有
 {"retcode":116,"p":"39bd5c71be123aaf451073154d52bfef78b1adaa0d087601"} // 什么都没有
 {
 	"retcode": 0,
 	"result": [{
 		"poll_type": "message",
 		"value": {...}
 	},{
 		"poll_type": "group_message",
 		"value": {...}
 	},{
 		"poll_type": "buddies_status_change",
 		"value": {...}
 	}]
 }
 poll2部分
 -----------
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

 "poll_type": "buddies_status_change"
 {
 	"uin": 3255435951,
 	"status": "online",
 	"client_type": 1
 }


 "poll_type": "group_web_message"
 {
	"msg_id": 27901,
	"from_uin": 2825512944,
	"to_uin": 2735284921,
	"msg_id2": 79523,
	"msg_type": 45,
	"reply_ip": 176756625,
	"group_code": 3447839986,
	"group_type": 1,
	"ver": 1,
	"send_uin": 630101664,
	"xml": "\u003c?xml version=\"1.0\" encoding=\"utf-8\"?\u003e\u003cd\u003e\u003cn t=\"h\" u=\"756458112\" i=\"103\" s=\"qun.qq.com/god/images/share-s.gif\" c=\"75797A\"/\u003e\u003cn t=\"t\" s=\"\u5171\u4EAB\u6587\u4EF6\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"t\" l=\"share\" s=\"1\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"t\" s=\"\u4E2A\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"t\" s=\"13:53:20\"/\u003e\u003cn t=\"r\"/\u003e\u003cn t=\"i\" s=\"qun.qq.com/god/images/share-other.gif\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"t\" s=\"cppsp_0.2.4.tar.xz\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"c\"/\u003e\u003cn t=\"r\"/\u003e\u003cn t=\"t\" l=\"share/download/!104!fa6db80f-86f4-4f40-bbcb-61729dbf1351$cppsp_0.2.4.tar.xz\" s=\"\u7ACB\u5373\u4E0B\u8F7D\"/\u003e\u003cn t=\"s\"/\u003e\u003cn t=\"t\" l=\"share\" s=\"\u67E5\u770B\u5168\u90E8\"/\u003e\u003cn t=\"r\"/\u003e\"/\u003e\u003c/d\u003e"
}
{
	"msg_id": 27933,
	"from_uin": 2825512944,
	"to_uin": 2735284921,
	"msg_id2": 79779,
	"msg_type": 45,
	"reply_ip": 176756625,
	"group_code": 3447839986,
	"group_type": 1,
	"ver": 1,
	"send_uin": 630101664,
	"xml": "\u003c?xml version=\"1.0\" encoding=\"utf-8\"?\u003e\u003cd\u003e\u003cn t=\"h\" u=\"756458112\" i=\"6\" s=\"1.url.cn/qun/feeds/img/server/g16.png\"/\u003e\u003cn t=\"t\" s=\"\u5206\u4EAB\u6587\u4EF6\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"t\" s=\"cppsp_0.2.4.tar.xz\"/\u003e\u003c/d\u003e"
}
{
	"msg_id": 27902,
	"from_uin": 2825512944,
	"to_uin": 2735284921,
	"msg_id2": 600924,
	"msg_type": 45,
	"reply_ip": 176488596,
	"group_code": 3447839986,
	"group_type": 1,
	"ver": 1,
	"send_uin": 630101664,
	"xml": "\u003c?xml version=\"1.0\" encoding=\"utf-8\"?\u003e\u003cd\u003e\u003cn t=\"h\" u=\"756458112\"/\u003e\u003cn t=\"t\" s=\"\u5171\u4EAB\u6587\u4EF6\"/\u003e\u003cn t=\"b\"/\u003e\u003cn t=\"t\" s=\"cppsp_0.2.4.tar.xz\"/\u003e\u003c/d\u003e"
}








*/
package webqq
