package target

import (
	"encoding/json"
	"github.com/pa001024/MoeWorker/source"
	"github.com/pa001024/MoeWorker/util"
)

type ITarget interface {
	Send(src *source.FeedInfo) (rid string, err error)
	GetMethod() []*TargetMethod
	GetId() string
}

type Target struct { // 配置持久模板
	Type   string          `json:"type"`   // 类型 target工厂ID 如[sinaweibo,qqweibo,renren]等
	Name   string          `json:"name"`   // 名字 不能跟别的target或source名字相同
	Method []*TargetMethod `json:"method"` // 处理方法
}
type TargetMethod struct {
	Action string   `json:"action"` // 动作 可选[update,repost]
	Source []string `json:"source"` // 目标 填写相应source或target名字
	Filter []string `json:"filter"` // 过滤器 从左到右依次管道
}

func New(name string, b []byte) (rst ITarget) {
	obj := &Target{}
	err := json.Unmarshal(b, obj)
	if err != nil {
		util.Log("JSON Parse Error", err)
		return
	}
	switch obj.Type {
	default:
	case "sina", "weibo", "sinaweibo":
		dst := &SinaWeibo{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.Log("target.sina \"" + name + "\" Loaded.")
	case "qqweibo", "qq":
		dst := &QQWeibo{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.Log("target.qqweibo \"" + name + "\" Loaded.")
	}
	return
}
