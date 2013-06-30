package filter

import (
	"github.com/pa001024/MoeCron/source"
)

type IFilter interface {
	Process(src []*source.FeedInfo) (dst []*source.FeedInfo)
}

type Filter struct {
	Type string `json:"type"` // 类型 filter工厂ID 如[moegirlwiki,简繁转换]等
	Name string `json:"name"` // 名字 不能跟别的target或source名字相同

	Impl IFilter `json:"-"` // 具体实现
}
