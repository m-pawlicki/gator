package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("Error forming requesst.")
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response.")
		os.Exit(1)
	}
	xmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body.")
		os.Exit(1)
	}
	rss := &RSSFeed{}
	err = xml.Unmarshal(xmlBytes, &rss)
	if err != nil {
		fmt.Println("Error unmarshalling XML.")
		os.Exit(1)
	}
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for key, val := range rss.Channel.Item {
		rss.Channel.Item[key].Title = html.UnescapeString(val.Title)
		rss.Channel.Item[key].Description = html.UnescapeString(val.Description)
	}
	return rss, nil
}
