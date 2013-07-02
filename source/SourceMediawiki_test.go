package source

import (
	"encoding/json"
	"log"
	"testing"
	"time"
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
			"interval":120,
			"limit":1,
			"pic":[100,100,800,600]
		}
		`), sour)
}

func TestGetFeed(t *testing.T) {
	// return
	sour.LastUpdate = time.Now().Add(-5 * time.Minute)
	f := sour.Get()
	for _, v := range f {
		log.Printf("%#v\n", v)
	}
	<-time.After(time.Second * 60)
	f2 := sour.Get()
	for _, v := range f2 {
		log.Printf("%#v\n", v)
	}
	if f[0].Content == f2[0].Content {
		t.Fatal("抓取重复!")
	}
}

func TestGetByName(t *testing.T) {
	// return
	log.Println(sour.GetByName("噬血狂襲"))
	log.Println(sour.GetByName("噬血狂襲22")) // test for not exists
}
