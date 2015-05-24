package tqq

import (
	"fmt"
)

type RemoteError string

type QQWeiboResult struct {
	ErrorCode int            `json:"errcode"` // 错误代码
	Error     string         `json:"msg"`     // 返回信息
	Ret       int            `json:"ret"`     // 返回值
	Data      *QQWeiboStatus `json:"data"`    // 数据
	// SeqId     string         `json:"seqid"`   // 序列号 (无需使用)
}
type QQWeiboStatus struct {
	Id        int64  `json:"id"`        // 微博id
	CreatedAt string `json:"timestamp"` // 微博发表时间
}

func (this RemoteError) Error() string {
	return "Remote Error: " + string(this)
}

// 转换到网页地址
func (this *QQWeiboStatus) Url() (urlText string) {
	urlText = fmt.Sprintf("http://t.qq.com/p/t/%v", this.Id)
	return
}
