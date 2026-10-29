package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadPost(t *testing.T) {
	source := `---
title: "Example Article"
date: "2026-10-30"
summary: "A short article summary."
---

## Introduction

This is the article content.
`

	path := writeFixture(t, t.TempDir(), "example.md", source)

	post, err := readPost(path)
	if err != nil {
		t.Fatalf("readPost() error = %v", err)
	}

	if post.Title != "Example Article" {
		t.Errorf("Title = %q", post.Title)
	}

	if post.Date != "2026-10-30" {
		t.Errorf("Date = %q", post.Date)
	}

	if post.Slug != "example" {
		t.Errorf("Slug = %q", post.Slug)
	}

	if !strings.Contains(string(post.Content), "<h2>Introduction</h2>") {
		t.Error("article heading was not converted to HTML")
	}
}

func TestReadPostRejectsInvalidMetadata(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "missing opening delimiter",
			source: "title: Example\n",
		},
		{
			name: "missing closing delimiter",
			source: `---
title: "Example"
date: "2026-10-30"
summary: "Summary"
`,
		},
		{
			name: "missing title",
			source: `---
date: "2026-10-30"
summary: "Summary"
---
Content
`,
		},
		{
			name: "missing date",
			source: `---
title: "Example"
summary: "Summary"
---
Content
`,
		},
		{
			name: "missing summary",
			source: `---
title: "Example"
date: "2026-10-30"
---
Content
`,
		},
		{
			name: "invalid date",
			source: `---
title: "Example"
date: "not-a-date"
summary: "Summary"
---
Content
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := writeFixture(t, t.TempDir(), "post.md", test.source)

			if _, err := readPost(path); err == nil {
				t.Fatal("readPost() expected an error")
			}
		})
	}
}

func TestLoadPostsSortsNewestFirst(t *testing.T) {
	directory := t.TempDir()

	writeFixture(t, directory, "older.md", `---
title: "Older Article"
date: "2026-10-10"
summary: "An older article."
---
Older content.
`)

	writeFixture(t, directory, "newer.md", `---
title: "Newer Article"
date: "2026-10-29"
summary: "A newer article."
---
Newer content.
`)

	posts, err := LoadPosts(directory)
	if err != nil {
		t.Fatalf("LoadPosts() error = %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("got %d posts, want 2", len(posts))
	}

	if posts[0].Slug != "newer" || posts[1].Slug != "older" {
		t.Errorf(
			"unexpected ordering: %q, %q",
			posts[0].Slug,
			posts[1].Slug,
		)
	}
}

func TestRenderMarkdownDoesNotRenderRawScriptTags(t *testing.T) {
	source := []byte("<script>alert(1)</script>")

	rendered, err := RenderMarkdown(source)
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	if strings.Contains(strings.ToLower(string(rendered)), "<script") {
		t.Errorf("raw script tag was rendered: %s", rendered)
	}
}

func TestRenderMarkdownRejectsJavaScriptLink(t *testing.T) {
	source := []byte("[click](javascript:alert%281%29)")

	rendered, err := RenderMarkdown(source)
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	output := strings.ToLower(string(rendered))

	if strings.Contains(output, `href="javascript:`) {
		t.Errorf("unsafe URL was rendered: %s", rendered)
	}
}

func writeFixture(t *testing.T, directory, name, source string) string {
	t.Helper()

	path := filepath.Join(directory, name)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	return path
}
