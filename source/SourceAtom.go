package source

import (
	"encoding/xml"
	"time"
)

// "format":"atom"
type FeedAtom struct {
	XMLName     xml.Name         `xml:"feed"`
	Id          string           `xml:"id"`
	Title       string           `xml:"title"`
	Updated     string           `xml:"updated"`
	Description string           `xml:"subtitle"`
	Generator   string           `xml:"generator"`
	Entry       []*FeedAtomEntry `xml:"entry"`
}
type FeedAtomEntry struct {
	Id            string `xml:"id"`
	Title         string `xml:"title"`
	AlternateLink string `xml:"link>alternate"`
	Updated       string `xml:"updated"`
	Summary       string `xml:"summary"`
	Author        string `xml:"author>name"`
}

type SourceAtom struct { // Atom 实现接口ISource
	// ISource
	Source

	FeedUrl    string    `json:"feed_url"` // http://www.mediawiki.org/wiki/Special:RecentChanges?feed=rss&namespace=0
	LastUpdate time.Time `json:"lastupdate"`
}
