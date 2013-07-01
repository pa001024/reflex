package filter

import (
	"github.com/pa001024/MoeCron/source"
)

type FilterBasic struct { // 基础过滤
	IFilter
	Filter

	MaxLength int    `json:"max_length"`
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
}

func (this *FilterBasic) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	dst = make([]*source.FeedInfo, len(src))
	for i, v := range src {
		nv := *v
		nv.Content = this.FilterContent(nv.Content)
		dst[i] = &nv
	}
	return
}

func (this *FilterBasic) FilterContent(src string) (dst string) {
	ds := []rune(dst)
	if len(ds) > this.MaxLength {
		dst = this.Prefix + string(ds[0:this.MaxLength]) + this.Suffix
	}
	return
}
