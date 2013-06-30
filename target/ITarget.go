package target

import (
	"github.com/pa001024/MoeCron/source"
)

type ITarget interface {
	Send(src *source.FeedInfo) bool
}

type Target struct {
	Type   string          `json:"type"`   // 类型 target工厂ID 如[sina,qq,renren]等
	Name   string          `json:"name"`   // 名字 不能跟别的target或source名字相同
	Method []*TargetMethod `json:"method"` // 处理方法

	Impl ITarget `json:"-"` // 具体实现
}
type TargetMethod struct {
	Action string   `json:"action"` // 动作 可选[update,upload]
	Source []string `json:"source"` // 目标 填写相应source或target名字
	Filter []string `json:"filter"` // 过滤器 从左到右依次管道
}
