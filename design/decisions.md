# Website Design Decisions

## 1. How the website runs

I decided to generate the website with Go and publish the generated
files as a static website.

The Go program reads the Markdown content, processes the templates,
and writes the finished pages to the `dist/` directory. The published
website consists of HTML, CSS, and supporting files.

For this project, a separate application server is not needed to
display the homepage, portfolio pages, or blog articles.

## 2. Why I chose static generation

The main purpose of this website is to present my work, share articles,
and provide a simple way for people to contact me.

Most of its content changes when I publish an update, not while a
visitor is viewing a page. Generating the pages ahead of time is
therefore a reasonable fit.

The main advantages are:

- The website can be hosted on a static hosting service.
- Visitors receive ready-made HTML pages.
- There is no application server to maintain in production.
- The build can validate content and generate the blog, RSS feed,
sitemap, and other files before deployment.
- The project has fewer moving parts than a dynamic web application.

## 3. Static generation versus a runtime server

A static generator creates the pages before visitors request them.
The web host then serves the generated files.

A runtime server generates a response when a request arrives. It is
useful when the application needs features such as authentication,
personalized responses, database queries, or other request-time logic.

For this website, most pages contain the same content for every visitor.
A runtime server would introduce additional deployment and maintenance
requirements without providing much benefit to the current scope.

| Area | Static generation | Runtime server |
| --- | --- | --- |
| Page generation | During the build | Usually when requested |
| Hosting | Static hosting is sufficient | Requires a server or suitable runtime |
| Content updates | Rebuild and redeploy | Can use updated data at runtime |
| Request-time logic | Limited without another service | Can handle application logic |
| Operational work | Relatively low | More server configuration and monitoring |
| Best fit here | Portfolio and blog | Features requiring live server-side behaviour |

## 4. Local development and production

During development, the Go generator creates the website in `dist/`.

The preview program serves that directory locally so I can inspect the
pages in a browser.

In production, GitHub Pages can host the generated files. The Go
generator runs during the build process; it does not run as a live
Go server on GitHub Pages.

A content change must be followed by another build and deployment
before the published website reflects that change.

## 5. What this means for the blog

Each article is stored in a Markdown file. The generator converts it
into an HTML page and adds it to the blog listing.

The same build produces the RSS feed and sitemap from the available
articles and pages.

This keeps the publishing process simple: edit the source content,
build the website, verify the output, and deploy it.

## 6. Trade-offs

Static generation does not automatically provide features such as
user accounts, server-side forms, or personalized content.

If I need those features later, I can evaluate a separate service or
a runtime backend for that specific requirement.

I do not need to add those components before there is a real use case
for them.

## 7. When to reconsider this decision

I will reconsider the architecture if the website needs features that
cannot reasonably be handled by generated pages and static hosting.

Examples include authenticated areas, user-submitted content that must
be processed securely, or content that must be fetched dynamically.

Until then, keeping the site static reduces operational complexity
and lets me focus on the content and quality of the website.