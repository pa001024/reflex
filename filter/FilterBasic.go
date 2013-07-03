package filter

import (
	"bytes"
	"github.com/pa001024/MoeCron/source"
	"io/ioutil"
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
		wr := &bytes.Buffer{}
		this.compFormat.Execute(wr, nv)
		s, _ := ioutil.ReadAll(wr)
		nv.Content = this.FilterContent(string(s))
		dst[i] = &nv
	}
	return
}

func (this *FilterBasic) FilterContent(src string) (dst string) {
	ds := []rune(dst)
	if len(ds) > this.MaxLength {
		dst = string(ds[0:this.MaxLength])
	}
	return
}
