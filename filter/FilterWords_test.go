package filter

import (
	"bytes"
	"fmt"
	"github.com/pa001024/MoeWorker/source"
	"testing"
)

func TestFilterWords(t *testing.T) {
	f := new(FilterWords)
	r := bytes.NewBufferString(`没有节操
节操掉啦
节操~`)
	f.LoadDict(r)

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
