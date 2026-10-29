# Website Architecture

## What this project does

This is my personal website and blog. I use Go to generate the website
from HTML templates and Markdown files. The generated files are saved
in `dist/` and can be hosted as a static website.

The homepage is kept simple. The other pages share a common layout,
and blog posts are generated from Markdown files.

## Project structure

```text
personal-site/
├── cmd/
│   ├── sitegen/
│   └── preview/
├── internal/
│   ├── build/
│   ├── config/
│   ├── content/
│   └── render/
├── content/
│   ├── pages/
│   └── posts/
├── templates/
├── static/
│   └── css/
├── design/
└── dist/