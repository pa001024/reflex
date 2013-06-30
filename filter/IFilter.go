package filter

import (
	"encoding/json"
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
)

type IFilter interface {
	Process(src []*source.FeedInfo) (dst []*source.FeedInfo)
}

type Filter struct {
	Type string `json:"type"` // 类型 filter工厂ID 如[moegirlwiki,简繁转换]等
	Name string `json:"name"` // 名字 不能跟别的target或source名字相同
}

func New(name string, b []byte) (rst IFilter) {
	obj := &Filter{}
	err := json.Unmarshal(b, obj)
	if err != nil {
		util.Log("JSON Parse Error", err)
		return
	}
	switch obj.Type {
	default:
	case "moegirlwiki":
		rst = &FilterMoegirlwiki{}
		json.Unmarshal(b, rst)
		rst.(*FilterMoegirlwiki).Name = name
		break
		// case "rss":
	}
	return
}
