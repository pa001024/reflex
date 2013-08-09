package filter

import (
	"github.com/pa001024/MoeWorker/source"
	"github.com/pa001024/MoeWorker/util"
	"regexp"
)

type FilterRegexp struct {
	IFilter
	Filter

	RegexpText string         `json:"regexp"`
	compFilter *regexp.Regexp `json:"-"`
}

func (this *FilterRegexp) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	if this.compFilter == nil {
		reg, err := regexp.Compile(this.RegexpText)
		if err != nil {
			util.ERROR.Err("Parse Regexp Fail", err)
			return src
		}
		this.compFilter = reg
	}
	dst = make([]*source.FeedInfo, 0, len(src))
	for _, v := range src {
		if !this.compFilter.MatchString(v.Title) && !this.compFilter.MatchString(v.Content) { // 需要优化
			dst = append(dst, v)
		}
	}
	return
}
