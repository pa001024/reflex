package filter

import (
	"github.com/pa001024/MoeWorker/source"
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
	dst = make([]*source.FeedInfo, len(src))
	for i, v := range src {
		nv := *v
		nv.Title = fss.Replace(nv.Title)
		nv.Content = fss.Replace(nv.Content)
		dst[i] = &nv
	}
	return
}
