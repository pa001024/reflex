package source

import (
	"encoding/xml"
	"time"
)

// "format":"rss"
type FeedRSS struct {
	XMLName     xml.Name       `xml:"rss>channel"`
	Id          string         `xml:"id"`
	Title       string         `xml:"title"`
	Updated     string         `xml:"lastBuildDate"`
	Description string         `xml:"description"`
	Generator   string         `xml:"generator"`
	Language    string         `xml:"language"`
	Item        []*FeedRSSItem `xml:"item"`
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
	LastUpdate time.Time `json:"lastupdate"`
}
