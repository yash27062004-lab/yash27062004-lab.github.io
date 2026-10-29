package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/config"
)

func TestBuildGeneratesWebsite(t *testing.T) {
	root := t.TempDir()

	writeFixture(t, root, "templates/layout.html", `{{define "layout"}}<!doctype html>
<html>
<head>
<title>{{.MetaTitle}}</title>
<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
<header>
<a href="/">Yash Sharma</a>
<nav>
<a href="/about/">About</a>
<a href="/expertise/">Expertise</a>
<a href="/experience/">Experience</a>
<a href="/blog/">Blog</a>
<a href="/contact/">Contact</a>
</nav>
</header>
<main>{{template "content" .}}</main>
</body>
</html>{{end}}`)

	writeFixture(t, root, "templates/home.html", `{{define "content"}}
<section class="home"><h1>Yash Sharma</h1><p>Software Engineer</p></section>
{{end}}`)

	writeFixture(t, root, "templates/page.html", `{{define "content"}}
<section><h1>{{.Title}}</h1><p>{{.Description}}</p><div>{{.Content}}</div></section>
{{end}}`)

	writeFixture(t, root, "templates/blog.html", `{{define "content"}}
<section><h1>{{.Title}}</h1>
{{range .Posts}}<article><h2><a href="/blog/{{.Slug}}/">{{.Title}}</a></h2><p>{{.Summary}}</p></article>{{end}}
</section>
{{end}}`)

	writeFixture(t, root, "templates/post.html", `{{define "content"}}
<article><h1>{{.Post.Title}}</h1><div>{{.Post.Content}}</div></article>
{{end}}`)

	for _, name := range []string{
		"about",
		"expertise",
		"experience",
		"contact",
		"not-found",
	} {
		writeFixture(
			t,
			root,
			filepath.Join("content", "pages", name+".md"),
			"Content for "+name+".\n",
		)
	}

	writeFixture(t, root, "content/posts/first-post.md", `---
title: "First Article"
date: "2026-10-30"
summary: "A test article."
---

## Article Body

Hello from the article.
`)

	writeFixture(
		t,
		root,
		"static/css/style.css",
		"body { color: black; }\n",
	)

	cfg := config.Config{
		OutputDir:    filepath.Join(root, "dist"),
		TemplatesDir: filepath.Join(root, "templates"),
		ContentDir:   filepath.Join(root, "content"),
		StaticDir:    filepath.Join(root, "static"),
		SiteURL:      "https://example.test",
	}

	if err := New(cfg).Build(); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	expectedFiles := []string{
		"index.html",
		"about/index.html",
		"expertise/index.html",
		"experience/index.html",
		"contact/index.html",
		"blog/index.html",
		"blog/first-post/index.html",
		"404.html",
		"feed.xml",
		"sitemap.xml",
		"robots.txt",
		"static/css/style.css",
	}

	for _, name := range expectedFiles {
		path := filepath.Join(cfg.OutputDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing generated file %q: %v", name, err)
		}
	}

	assertContains(t, filepath.Join(cfg.OutputDir, "index.html"), `href="/about/"`)
	assertContains(t, filepath.Join(cfg.OutputDir, "index.html"), `href="/static/css/style.css"`)
	assertContains(t, filepath.Join(cfg.OutputDir, "blog/index.html"), "First Article")
	assertContains(t, filepath.Join(cfg.OutputDir, "blog/first-post/index.html"), "Article Body")
	assertContains(t, filepath.Join(cfg.OutputDir, "feed.xml"), "/blog/first-post/")
	assertContains(t, filepath.Join(cfg.OutputDir, "sitemap.xml"), "/blog/first-post/")
	assertContains(t, filepath.Join(cfg.OutputDir, "robots.txt"), "/sitemap.xml")
}

func TestBuildRejectsInvalidConfiguration(t *testing.T) {
	err := New(config.Config{}).Build()

	if err == nil {
		t.Fatal("Build() expected an error for invalid configuration")
	}

	if !strings.Contains(err.Error(), "validate configuration") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuildProtectsSourceDirectories(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "content")

	writeFixture(
		t,
		root,
		filepath.Join("content", "pages", "about.md"),
		"Do not delete this file.\n",
	)

	cfg := config.Config{
		OutputDir:    contentDir,
		TemplatesDir: filepath.Join(root, "templates"),
		ContentDir:   contentDir,
		StaticDir:    filepath.Join(root, "static"),
		SiteURL:      "https://example.test",
	}

	err := New(cfg).Build()
	if err == nil {
		t.Fatal("Build() expected an error for overlapping paths")
	}

	if !strings.Contains(err.Error(), "overlaps source directory") {
		t.Errorf("unexpected error: %v", err)
	}

	sourcePath := filepath.Join(contentDir, "pages", "about.md")
	if _, err := os.Stat(sourcePath); err != nil {
		t.Fatalf("source file was unexpectedly removed: %v", err)
	}
}

func writeFixture(t *testing.T, root, name, source string) {
	t.Helper()

	path := filepath.Join(root, name)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}

func assertContains(t *testing.T, path, expected string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}

	if !strings.Contains(string(data), expected) {
		t.Errorf("%q does not contain %q", path, expected)
	}
}
