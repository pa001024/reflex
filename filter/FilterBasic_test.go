package filter

import (
	"fmt"
	"github.com/pa001024/reflex/source"
	"testing"
)

const (
	ex_f = `
		{
			"type":"basic",
			"max_length":100,
			"format":"{{.Content}}... 阅读更多: {{urlquery .Link}}"
		}
		`
)

func TestBasic(t *testing.T) {
	f := New("name", []byte(ex_f))
	src := []*source.FeedInfo{
		&source.FeedInfo{
			Link:    "http://wwww",
			Content: "内容简介之类的说...那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥那啥长度超了的话",
		},
	}
	ns := f.Process(src)
	fmt.Printf("%#v\n", src[0])
	fmt.Printf("%#v\n", ns[0])
}
