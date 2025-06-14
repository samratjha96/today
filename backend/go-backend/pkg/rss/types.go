package rss

// RSSItem represents a news item from an RSS feed
type RSSItem struct {
	Source string `json:"source"`
	Title  string `json:"title"`
	Link   string `json:"link"`
}

// RSSFeed represents a news source with its URL
type RSSFeed struct {
	Name string
	URL  string
}

// RSSFeedList is a map of feed names to feed URLs
var RSSFeedList = map[string]string{
	"TechCrunch":   "http://feeds.feedburner.com/TechCrunch/",
	"Wired":        "https://www.wired.com/feed/rss",
	"The Verge":    "https://www.theverge.com/rss/index.xml",
	"Ars Technica": "http://feeds.arstechnica.com/arstechnica/index",
}
