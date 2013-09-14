package filter

import (
	"fmt"
	"github.com/pa001024/reflex/source"
	"testing"
)

func TestFilterRegexp(t *testing.T) {
	f := new(FilterRegexp)
	f.RegexpText = `节操~?|哟{4,}`

	src := []*source.FeedInfo{
		&source.FeedInfo{
			Title:   "没有节操",
			Content: `诶嘿嘿~~~`,
		},
		&source.FeedInfo{
			Title:   "哟",
			Content: `节操~`,
		},
		&source.FeedInfo{
			Title:   "哟",
			Content: `哟哟哟~`,
		},
		&source.FeedInfo{
			Title:   "哟哟哟",
			Content: `哟哟哟哟~`,
		},
	}
	final := f.Process(src)
	for _, v := range final {
		fmt.Printf("%#v\n", v)
	}

}
