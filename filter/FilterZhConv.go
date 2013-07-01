package filter

import (
	"github.com/pa001024/MoeCron/source"
	"strings"
)

type FilterZhConv struct { // 简繁转换
	IFilter
	Filter
}

var (
	fss *strings.Replacer
)

func (this *FilterZhConv) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	if fss == nil {
		fss = strings.NewReplacer(zh_t2s_table...)
	}
	dst = make([]*source.FeedInfo, 0, len(src))
	for i, v := range src {
		nv := *v
		nv.Content = fss.Replace(nv.Content)
		dst[i] = &nv
	}
	return
}
