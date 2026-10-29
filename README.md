# Personal Website

A personal website and technical blog generated with Go from Markdown content and HTML templates.

**Live website:** https://yash27062004-lab.github.io
**Source repository:** https://github.com/yash27062004-lab/yash27062004-lab.github.io

## Features

- Responsive personal homepage and content pages
- Markdown-based blog posts
- Go-powered static site generation
- RSS feed, sitemap, and `robots.txt`
- Automated formatting, static analysis, tests, build checks, and GitHub Pages deployment

## Project structure

```text
cmd/
preview/       Local preview server
sitegen/       Static site generator entry point
content/
pages/         Website page content
posts/         Markdown blog posts
design/          Architecture, design decisions, and wireframes
internal/
build/         Site generation and feeds
config/        Build configuration
content/       Markdown and post processing
render/        Template rendering and static file copying
static/css/      Stylesheet
templates/       HTML templates
dist/            Generated website output (not committed)
```

## Requirements

Install a Go version compatible with the `go` version declared in `go.mod`.

## Build, test, and preview

From the repository root, run:

```powershell
go mod tidy
go run ./cmd/sitegen
go vet ./...
go test ./...
go build ./...
go run ./cmd/preview
```

The generator writes the static website to `dist/`. The preview server serves that directory locally; open `http://127.0.0.1:8080` after starting it. Run the generator first, and start the preview server in a terminal where you can see its logs.

## Writing a blog post

Create a Markdown file under `content/posts/`, using the file name as the URL slug. Each post needs `title`, `date` (`YYYY-MM-DD`), and `summary` fields between opening and closing `---` front-matter delimiters.

Example:

```markdown
---
title: My First Post
date: 2026-10-30
summary: A short description of the article.
---

# My First Post

Write the article here.
```

After saving the file, run `go run ./cmd/sitegen` and inspect the generated blog listing, article page, RSS feed, and sitemap before committing.

## Configuration

The generator supports these environment-variable overrides:

- `SITE_URL`
- `SITE_OUTPUT_DIR`
- `SITE_TEMPLATES_DIR`
- `SITE_CONTENT_DIR`
- `SITE_STATIC_DIR`

Defaults are defined in `internal/config/config.go`.

## Deployment

The GitHub Actions workflow on the `main` branch validates the Go code, generates the site, uploads the `dist/` directory, and deploys it to GitHub Pages. In repository settings, GitHub Pages should use **GitHub Actions** as its deployment source.

Push changes to `main` only after the local tests and build pass. Review the Actions run to confirm the deployment succeeded.

## Licensing

- **Original software code** (including the Go generator, HTML templates, stylesheet, and build workflow) is licensed under the MIT License in `LICENSE`.
- **Original editorial content displayed on the website**—including blog articles, page copy, biography text, and original creative material—is covered by `CONTENT-LICENSE.md` and is all rights reserved unless a specific license says otherwise.
- **Third-party components** retain their own license terms; see `THIRD_PARTY_NOTICES.md`.

The above distinctions describe the intended scope of each license; this is not legal advice. The licenses do not change the separate terms, limits, or pricing policies of hosting and other services.