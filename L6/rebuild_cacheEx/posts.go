package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL      string   `xml:"guid"`
	Title    string   `xml:"title"`
	Category []string `xml:"category"`
}

func FetchContent(url string) (*RSS, error) {
	fmt.Printf("fetching: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET for URL %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}
	return rss, nil
}
