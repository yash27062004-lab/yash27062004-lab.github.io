package build

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/content"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	PubDate     string  `xml:"pubDate"`
	GUID        RSSGUID `xml:"guid"`
}

type RSSGUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type Sitemap struct {
	XMLName xml.Name       `xml:"urlset"`
	XMLNS   string         `xml:"xmlns,attr"`
	URLs    []SitemapEntry `xml:"url"`
}

type SitemapEntry struct {
	Location string `xml:"loc"`
}

func writeRSS(path, siteURL string, posts []content.Post) error {
	feed := RSSFeed{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       "Yash Sharma",
			Link:        siteURL,
			Description: "Articles and writing by Yash Sharma.",
			Language:    "en",
			Items:       make([]RSSItem, 0, len(posts)),
		},
	}

	for _, post := range posts {
		date, err := time.Parse("2006-01-02", post.Date)
		if err != nil {
			return fmt.Errorf("invalid date for article %q: %w", post.Title, err)
		}

		postURL := siteURL + "/blog/" + post.Slug + "/"

		feed.Channel.Items = append(feed.Channel.Items, RSSItem{
			Title:       post.Title,
			Link:        postURL,
			Description: post.Summary,
			PubDate:     date.UTC().Format(time.RFC1123Z),
			GUID: RSSGUID{
				IsPermaLink: true,
				Value:       postURL,
			},
		})
	}

	return writeXML(path, feed)
}

func writeSitemap(path, siteURL string, paths []string) error {
	sort.Strings(paths)

	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapEntry, 0, len(paths)),
	}

	for _, route := range paths {
		sitemap.URLs = append(sitemap.URLs, SitemapEntry{
			Location: siteURL + route,
		})
	}

	return writeXML(path, sitemap)
}

func writeRobots(path, siteURL string) error {
	robotsContent := "User-agent: *\n" +
		"Allow: /\n" +
		"Sitemap: " + siteURL + "/sitemap.xml\n"

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create robots directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(robotsContent), 0644); err != nil {
		return fmt.Errorf("write robots.txt: %w", err)
	}

	return nil
}

func writeXML(path string, value any) error {
	data, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal XML: %w", err)
	}

	data = append([]byte(xml.Header), data...)
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write XML file: %w", err)
	}

	return nil
}
