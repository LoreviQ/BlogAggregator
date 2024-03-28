package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Item        []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (cfg apiConfig) startScraper() error {
	log.Printf("Scraping %v feeds every %v\n", cfg.noFeeds, cfg.interval.String())
	ticker := time.NewTicker(cfg.interval)
	for ; ; <-ticker.C {
		feeds, err := cfg.DB.GetNextFeeds(context.Background(), cfg.noFeeds)
		if err != nil {
			return err
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go cfg.scrapeFeed(wg, feed)
		}
		wg.Wait()
	}
}

func (cfg apiConfig) scrapeFeed(wg *sync.WaitGroup, feedDB database.Feed) {
	defer wg.Done()
	feedStruct, err := getFeedFromEndpoint(feedDB.Url)
	cfg.DB.UpdateFetched(context.Background(), feedDB.ID)
	if err != nil {
		log.Printf("Failed to obtain feed from %v - Err: %v\n", feedDB.Url, err)
		return
	}
	for _, post := range feedStruct.Channel.Item {
		cfg.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       convertToNullString(post.Title),
			Url:         post.Link,
			Description: convertToNullString(post.Description),
			PublishedAt: parseTime(post.PubDate),
			FeedID:      feedDB.ID,
		})
	}
}

func getFeedFromEndpoint(endpoint string) (*RSSFeed, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	var rssFeed RSSFeed
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return nil, err
	}
	return &rssFeed, nil
}

func convertToNullString(s string) sql.NullString {
	nullString := sql.NullString{
		String: s,
	}
	if s == "" {
		nullString.Valid = false
	} else {
		nullString.Valid = true
	}
	return nullString
}

func parseTime(t string) sql.NullTime {
	if t == "" {
		return sql.NullTime{
			Valid: false,
		}
	}
	layouts := map[string]string{
		"layout RFC 822 ver1":  "02 Jan 06 15:04 MST",
		"layout RFC 822 ver2":  "02 Jan 06 15:04 -0700",
		"layout RFC 1123 ver1": "Mon, 02 Jan 2006 15:04:05 MST",
		"layout RFC 1123 ver2": "Mon, 02 Jan 2006 15:04:05 -0700",
		"layout RFC 1123 ver3": "Mon, 2 Jan 2006 15:04:05 MST",
		"layout RFC 1123 ver4": "Mon, 2 Jan 2006 15:04:05 -0700",
	}
	for _, layout := range layouts {
		parsedTime, err := time.Parse(layout, t)
		if err == nil {
			return sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}
	}
	log.Printf("Failed to parse time\ntime: %v\n", t)
	return sql.NullTime{
		Valid: false,
	}
}
