package filter

import (
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type FilterWords struct {
	IFilter
	Filter

	DictFile   string            `json:"dict_file"`
	compFilter *strings.Replacer `json:"-"`
}

func (this *FilterWords) LoadDict(r io.Reader) error {
	g := regexp.MustCompile(`(\S.+) +(\S.+)`)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	m := g.FindAllStringSubmatch(string(b), -1)
	a := make([]string, 0, len(m)*2)
	for _, v := range m {
		if len(v) > 2 && v[1] != "" && v[2] != "" {
			a = append(a, v[1])
			a = append(a, v[2])
		}
	}
	this.compFilter = strings.NewReplacer(a...)
	return nil
}
func (this *FilterWords) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	if this.compFilter == nil {
		r, err := os.Open(this.DictFile)
		if err != nil && this.LoadDict(r) != nil {
			util.Log("Warning: FilterWords.LoadDict Fail", err)
			return src
		}
	}
	dst = make([]*source.FeedInfo, 0, len(src))
	for _, v := range src {
		if this.compFilter.Replace(v.Title) == v.Title && this.compFilter.Replace(v.Content) == v.Content { // TODO: 需要优化性能
			dst = append(dst, v)
		}
	}
	return
}
