package util

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type oEmbedResponse struct {
	Title string `json:"title"`
}

// ExtractTitle extracts a human-readable title for a URL.
// Order:
// 1. YouTube oEmbed (for accurate video titles)
// 2. HTML <title> tag
// 3. Domain fallback
func ExtractTitle(link string) string {
	// Special handling for YouTube
	if isYouTube(link) {
		if t := youtubeTitle(link); t != "" {
			return t
		}
	}

	// Generic HTML title extraction
	if t := htmlTitle(link); t != "" {
		return t
	}

	return domainFallback(link)
}

// --------------------
// YouTube handling
// --------------------

func isYouTube(link string) bool {
	return strings.Contains(link, "youtube.com") || strings.Contains(link, "youtu.be")
}

func youtubeTitle(link string) string {
	oembedURL := "https://www.youtube.com/oembed?format=json&url=" + url.QueryEscape(link)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(oembedURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var data oEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}

	return strings.TrimSpace(data.Title)
}

// --------------------
// HTML title fallback
// --------------------

func htmlTitle(link string) string {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(link)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ""
	}

	var title string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if title != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return strings.TrimSpace(title)
}

// --------------------
// Domain fallback
// --------------------

func domainFallback(link string) string {
	trimmed := strings.TrimPrefix(link, "https://")
	trimmed = strings.TrimPrefix(trimmed, "http://")
	parts := strings.Split(trimmed, "/")
	return parts[0]
}
