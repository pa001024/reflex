package filter

import ()

type FilterT2S struct { // 简繁转换
	IFilter
	Filter
}

func (this *FilterT2S) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	// 填坑吧!
}
