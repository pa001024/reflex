package filter

import (
	"encoding/json"
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
	"text/template"
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
		dst := &FilterMoegirlwiki{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.Log("filter.moegirlwiki \"" + name + "\" Loaded.")
	case "zhconv":
		dst := &FilterZhConv{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.Log("filter.zhconv \"" + name + "\" Loaded.")
	case "url":
		dst := &FilterUrl{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.Log("filter.url \"" + name + "\" Loaded.")
	case "basic":
		dst := &FilterBasic{}
		json.Unmarshal(b, dst)
		dst.Name = name
		dst.compFormat = template.Must(template.New(name).Parse(dst.Format))
		if dst.MaxLength == 0 {
			dst.MaxLength = 120
		}
		rst = dst
		util.Log("filter.basic \"" + name + "\" Loaded.")
	}
	return
}
