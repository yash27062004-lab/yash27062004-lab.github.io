---
title: Building My Personal Website with Go
date: 2026-10-30
summary: Why I built my personal website with Go, how the site works, and some problems I faced while building and deploying it.
---

# Building My Personal Website with Go

## Why I wanted a personal website

I wanted to build my own website where I can share my work, write about things I learn and have a place to show my projects. A normal resume is useful, but it doesn't really explain how I work or what kind of problems I enjoy solving.

I could have used a website builder or downloaded a ready-made theme, but I wanted to build something on my own. I also wanted to understand what happens behind the scenes when a website is generated and published.

I decided to keep the design simple. I don't want too many animations or unnecessary things on the page. The main focus should be on my work, my writing and the information people come here to find.

## What I wanted to build

My goal was to make a personal website with a homepage, some content pages and a blog where I can publish technical articles.

I wanted the website to be easy to maintain. Writing a new blog post shouldn't require me to change the Go code every time. I should be able to write the content in a Markdown file, run the generator and get a new page.

I also wanted to host it publicly without paying for a server. Since this is mainly a personal website, I didn't need a backend server running all the time just to serve pages.

That is why I decided to build a static site generator using Go and host the generated website on GitHub Pages.

## How the site works

The basic idea is pretty simple. I write the content, Go processes it, and the generator creates the HTML files that will be published on the internet.

The project has different folders for different responsibilities.

The `content/pages/` directory contains the website pages, while `content/posts/` is where blog articles will live. Blog posts use Markdown, which makes writing easier because I can focus on the content instead of writing HTML for everything.

The `templates/` directory contains the HTML templates used to generate the pages. This lets me keep the general layout separate from the actual content.

The Go code reads the configuration and content, processes the Markdown, applies the templates and writes the final website into the `dist/` directory. Static files like CSS are copied there as well.

I also added RSS, a sitemap and a robots.txt file. These are small things, but they make the website more complete and easier to maintain.

Once the website is generated, GitHub Pages can serve those files directly. There is no need to run my Go program on a server every time someone visits the website.

## A problem I encountered

One thing I learned pretty quickly is that getting the code to work locally doesn't mean everything will work in the deployment pipeline.

At one point, `go build ./...` was failing because there was a root-level `main.go` file that only contained `package main`. It didn't have a `func main()`, so Go was treating it like an executable package without an entry point. I removed that unnecessary file because the actual commands already lived inside the `cmd/` directory.

Later, the GitHub Actions workflow failed while checking the generated website files. The Go build and tests were passing, so I looked at the failed step instead of changing random parts of the code.

The problem was a wrong file path in the workflow. My generator copied the stylesheet to `dist/static/css/style.css`, but the workflow was checking `dist/css/style.css`. The file was there, just not at the location the workflow expected.

I corrected the path and pushed the change again.

These issues were not very complicated, but they reminded me that debugging is often about checking the actual error and understanding what the program is doing. Sometimes the problem is not in the code you first suspect.

## What I learned

Building this website helped me understand how different parts of a small Go project fit together. It is one thing to learn individual Go concepts, but putting them into a real project makes you think more carefully about how files, packages and responsibilities should be organised.

I also got more practice with testing. I added tests for different parts of the project and started using commands like `go test ./...`, `go vet ./...` and `go build ./...` to check the code before pushing changes.

Another useful part was setting up GitHub Actions. The workflow checks the code, runs tests, builds the site and verifies that important generated files exist. If those checks pass, it uploads the generated website and deploys it to GitHub Pages.

It took a few attempts to get everything right, but now I have a process that can catch problems before the new version goes live.

I also learned that deployment is not just about uploading files. The generated paths, the configuration and the deployment workflow all need to agree with each other.

## What comes next

The website is now live, but I don't consider the project finished. There are still things I want to improve and learn along the way.

My next focus is to write more technical articles based on things I actually work on and understand. I also want to keep improving the tests, documentation and overall quality of the project as it grows.

I don't want to add features just to make the project look more complicated. I'd rather keep it simple and understand every part of what I'm building.

For now, this website is a starting point. I'll keep working on it, learn from the problems I run into, and share what I learn here.