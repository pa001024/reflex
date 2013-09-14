package filter

import (
	"bytes"
	"github.com/pa001024/reflex/source"
	"text/template"
)

type FilterBasic struct { // 基础过滤
	IFilter
	Filter

	MaxLength  int                `json:"max_length"`
	Format     string             `json:"format"`
	compFormat *template.Template `json:"-"`
}

func (this *FilterBasic) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	dst = make([]*source.FeedInfo, len(src))
	for i, v := range src {
		nv := *v
		nv.Content = this.FilterContent(nv.Content)
		buf := new(bytes.Buffer)
		this.compFormat.Execute(buf, nv)
		nv.Content = buf.String()
		dst[i] = &nv
	}
	return
}

func (this *FilterBasic) FilterContent(src string) (dst string) {
	if this.MaxLength <= 0 {
		return src
	}
	ds := []rune(src)
	if len(ds) > this.MaxLength {
		dst = string(ds[0:this.MaxLength])
	}
	return
}
