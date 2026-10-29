package content

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Slug    string
	Content template.HTML
}

func LoadPosts(directory string) ([]Post, error) {
	files, err := filepath.Glob(filepath.Join(directory, "*.md"))
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(files))

	for _, file := range files {
		post, err := readPost(file)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	return posts, nil
}

func ReadMarkdown(path string) (template.HTML, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return RenderMarkdown(source)
}

func RenderMarkdown(source []byte) (template.HTML, error) {
	var output bytes.Buffer

	if err := goldmark.Convert(source, &output); err != nil {
		return "", err
	}

	return template.HTML(output.String()), nil
}

func readPost(path string) (Post, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}

	source := strings.ReplaceAll(string(raw), "\r\n", "\n")
	lines := strings.Split(source, "\n")

	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return Post{}, fmt.Errorf("missing opening front matter delimiter")
	}

	var post Post
	metadataEnd := -1

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "---" {
			metadataEnd = i
			break
		}

		if line == "" {
			continue
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return Post{}, fmt.Errorf("invalid metadata line %q", line)
		}

		value = strings.Trim(strings.TrimSpace(value), "\"'")

		switch strings.TrimSpace(key) {
		case "title":
			post.Title = value
		case "date":
			post.Date = value
		case "summary":
			post.Summary = value
		}
	}

	if metadataEnd == -1 {
		return Post{}, fmt.Errorf("missing closing front matter delimiter")
	}

	if post.Title == "" || post.Date == "" || post.Summary == "" {
		return Post{}, fmt.Errorf("title, date, and summary are required")
	}

	if _, err := time.Parse("2006-01-02", post.Date); err != nil {
		return Post{}, fmt.Errorf("invalid date %q: %w", post.Date, err)
	}

	post.Slug = strings.TrimSuffix(
		filepath.Base(path),
		filepath.Ext(path),
	)

	body := []byte(strings.Join(lines[metadataEnd+1:], "\n"))

	post.Content, err = RenderMarkdown(body)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}
