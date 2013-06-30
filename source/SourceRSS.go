package source

import (
	"encoding/xml"
	"time"
)

// "format":"rss"
type FeedRSS struct {
	XMLName     xml.Name       `xml:"rss"`
	Id          string         `xml:"channel>id"`
	Title       string         `xml:"channel>title"`
	Updated     string         `xml:"channel>lastBuildDate"`
	Description string         `xml:"channel>description"`
	Generator   string         `xml:"channel>generator"`
	Language    string         `xml:"channel>language"`
	Item        []*FeedRSSItem `xml:"channel>item"`
}
type FeedRSSItem struct {
	Id          string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Updated     string `xml:"pubDate"`
	Description string `xml:"description"`
	Author      string `xml:"dc:creator"`
	Comments    string `xml:"comments"`
}

type SourceRSS struct { // RSS 实现接口ISource
	// ISource
	Source

	FeedUrl    string    `json:"feed_url"` // http://www.mediawiki.org/wiki/Special:RecentChanges?feed=rss&namespace=0
	LastUpdate time.Time `json:"-"`
}
