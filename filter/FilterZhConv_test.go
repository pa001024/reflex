package filter

import (
	"fmt"
	"github.com/pa001024/MoeCron/source"
	"testing"
)

const (
	ex_f3 = `
		{
			"type":"zhconv"
		}
		`
)

func TestZHConv(t *testing.T) {
	f := New("name", []byte(ex_f3))
	src := []*source.FeedInfo{
		&source.FeedInfo{
			Title:   "吳淞鐵路",
			Link:    "http://zh.wikipedia.org/wiki/%E5%90%B3%E6%B7%9E%E9%90%B5%E8%B7%AF",
			Content: `吳淞鐵路: 吳淞鐵路連接上海至吳淞鎮，是中國第一條辦理營業的鐵路。1876年由英國商人未經批准建造。通車十六個月後被清朝官員以二十八萬五千兩白銀買回，之後被拆除。 早在1825年，英國便建成世界上第一條投入商業營`,
		},
	}
	ns := f.Process(src)
	fmt.Printf("%#v\n", ns[0])
}
