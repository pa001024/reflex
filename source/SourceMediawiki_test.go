package source

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var (
	sour = &SourceMediawiki{}
)

func init() {
	json.Unmarshal([]byte(`
		{
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
	sour.LastUpdate = time.Now().Add(-5 * time.Minute)
	v := sour.Get()
	fmt.Printf("%+v\n", v)
}
