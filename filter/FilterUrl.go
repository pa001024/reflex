package filter

import (
	"github.com/pa001024/MoeCron/source"
	"net/url"
)

type FilterUrl struct { // 将URL进行encode 以修复微博等处的BUG
	IFilter
	Filter
}

func (this *FilterUrl) FilterLink(src string) (rst string) {
	return url.QueryEscape(src)
}
func (this *FilterUrl) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	dst = make([]*source.FeedInfo, len(src))
	for i, v := range src {
		nv := *v
		nv.Link = this.FilterLink(nv.Link)
		dst[i] = &nv
	}
	return
}
