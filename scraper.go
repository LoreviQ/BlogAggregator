package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
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
	log.Printf("Found the feed: %v\n", feedStruct.Channel.Title)
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
