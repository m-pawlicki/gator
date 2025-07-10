package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-pawlicki/gator/internal/database"
	"github.com/m-pawlicki/gator/internal/state"
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

func ScrapeFeeds(s *state.State) error {
	ctx := context.Background()
	nextFeed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Println("Failed to get next feed.")
		os.Exit(1)
	}
	markedFeed := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        nextFeed.ID,
	}
	err = s.DB.MarkFeedFetched(ctx, markedFeed)
	if err != nil {
		fmt.Println("Failed to mark feed as fetched.")
		os.Exit(1)
	}
	feedItems, err := FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		fmt.Println("Couldn't fetch feed.")
		os.Exit(1)
	}

	for _, item := range feedItems.Channel.Item {
		timeConv, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {

		}
		post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   nextFeed.CreatedAt,
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: timeConv,
			FeedID:      nextFeed.ID,
		}
		_, err = s.DB.CreatePost(ctx, post)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v\n", err)
			continue
		}
	}
	return nil
}
