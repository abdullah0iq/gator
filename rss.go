package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	rssFeed := &RSSFeed{}
	if err := xml.Unmarshal(data, rssFeed); err != nil {
		return &RSSFeed{}, err
	}
	return decodeEscapedHTML(rssFeed), nil
}
func decodeEscapedHTML(rssFeed *RSSFeed) (*RSSFeed) {
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i , rssItem := range rssFeed.Channel.Items {
		rssFeed.Channel.Items[i].Title =  html.UnescapeString(rssItem.Title)
		rssFeed.Channel.Items[i].Description =  html.UnescapeString(rssItem.Description)
	}
	return rssFeed
}
