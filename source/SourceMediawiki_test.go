package source

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

var (
	sour = &SourceMediawiki{}
)

func init() {
	json.Unmarshal([]byte(`
		{
			"name":"src1",
			"type":"mediawiki",
			"feed_url":"http://zh.wikipedia.org/w/index.php?title=Special:%E6%9C%80%E8%BF%91%E6%9B%B4%E6%94%B9&feed=rss&namespace=0",
			"api_url":"http://zh.wikipedia.org/w/api.php",
			"pic_url":"http://upload.wikimedia.org/wikipedia/zh/",
			"interval":120,
			"limit":1
		}
		`), sour)
}
func TestGet(t *testing.T) {
	// return
	lastid := ""
	for i := 0; i < 10; i++ {
		f := sour.Get()
		if len(f) > 0 {
			log.Print("")
		}
		for _, v := range f {
			if v.Id == lastid {
				t.Error("重复抓取!!!")
				t.Fail()
			}
			lastid = v.Id
			b := []rune(v.Content)
			if len(b) > 80 {
				v.Content = string(b[0:80])
			}
			fmt.Printf("%#v\n", v)
		}
	}
}

func TestGetByName(t *testing.T) {
	// return
	log.Println(sour.GetByName("噬血狂襲"))
	log.Println(sour.GetByName("噬血狂襲22")) // test for not exists
}

func TestFilterPicUrl(t *testing.T) {
	d := `[[image:ia asd.png]]
| image    = some asd as.png
http://www.com/jsp.jpg
`
	a := sour.FilterPicUrl(d)
	// fmt.Printf("%#v\n", d)
	fmt.Printf("%#v\n", a)
}
