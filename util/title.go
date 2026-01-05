package util

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func ExtractTitle(url string) string {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		return domainFallback(url)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return domainFallback(url)
	}

	var title string

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	if title == "" {
		return domainFallback(url)
	}

	return strings.TrimSpace(title)
}

func domainFallback(url string) string {
	parts := strings.Split(strings.TrimPrefix(url, "https://"), "/")
	return parts[0]
}
