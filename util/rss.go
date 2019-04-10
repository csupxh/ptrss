package util

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type RSS struct {
	xml.Name `xml:"rss"`
	Channel  RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	PubDate     string    `xml:"pubDate"`
	Generator   string    `xml:"generator"`
	Docs        string    `xml:"docs"`
	TTL         int       `xml:"ttl"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title           string       `xml:"title"`
	Link            string       `xml:"link"`
	Description     string       `xml:"description"`
	Enclosure       RSSEnclosure `xml:"enclosure"`
	GuidValue       string       `xml:"guid"`
	GuidIsPermaLink bool         `xml:"isPermaLink,attr"`
	PubDate         string       `xml:"pubDate"`
}

type RSSEnclosure struct {
	Url    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type ChannelCallback func(channel RSSChannel)

// 根据url获取xml转为RSS对象
func GetBody(url string, callback ChannelCallback) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var rss RSS
	errs := xml.Unmarshal(body, &rss)
	if errs != nil {
		log.Println(errs)
	}
	callback(rss.Channel)
}
