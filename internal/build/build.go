package build

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/config"
	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/content"
	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/render"
)

type PageData struct {
	Title       string
	MetaTitle   string
	Description string
	Email       string
	Content     template.HTML
}

type BlogData struct {
	PageData
	Posts []content.Post
}

type ArticleData struct {
	PageData
	Post content.Post
}

type pageDefinition struct {
	Slug        string
	Page        PageData
	ContentFile string
}

type Generator struct {
	cfg config.Config
}

func New(cfg config.Config) *Generator {
	cfg.SiteURL = strings.TrimRight(cfg.SiteURL, "/")
	return &Generator{cfg: cfg}
}

func (g *Generator) Build() error {
	if err := g.cfg.Validate(); err != nil {
		return fmt.Errorf("validate configuration: %w", err)
	}

	if err := g.prepareOutputDirectory(); err != nil {
		return fmt.Errorf("prepare output directory: %w", err)
	}

	slog.Info("generating homepage")

	if err := g.generateHome(); err != nil {
		return err
	}

	pages := g.pageDefinitions()

	if err := g.generatePages(pages); err != nil {
		return err
	}

	posts, err := content.LoadPosts(filepath.Join(g.cfg.ContentDir, "posts"))
	if err != nil {
		return fmt.Errorf("load blog posts: %w", err)
	}

	if err := g.generateBlog(posts); err != nil {
		return err
	}

	if err := g.generateArticles(posts); err != nil {
		return err
	}

	if err := g.generateNotFound(); err != nil {
		return err
	}

	if err := writeRSS(
		filepath.Join(g.cfg.OutputDir, "feed.xml"),
		g.cfg.SiteURL,
		posts,
	); err != nil {
		return fmt.Errorf("generate RSS feed: %w", err)
	}

	if err := writeSitemap(
		filepath.Join(g.cfg.OutputDir, "sitemap.xml"),
		g.cfg.SiteURL,
		g.sitemapPaths(pages, posts),
	); err != nil {
		return fmt.Errorf("generate sitemap: %w", err)
	}

	if err := writeRobots(
		filepath.Join(g.cfg.OutputDir, "robots.txt"),
		g.cfg.SiteURL,
	); err != nil {
		return fmt.Errorf("generate robots.txt: %w", err)
	}

	if err := render.CopyDir(
		g.cfg.StaticDir,
		filepath.Join(g.cfg.OutputDir, "static"),
	); err != nil {
		return fmt.Errorf("copy static assets: %w", err)
	}

	slog.Info(
		"website build completed",
		"output_dir", g.cfg.OutputDir,
		"articles", len(posts),
	)

	return nil
}

func (g *Generator) prepareOutputDirectory() error {
	outputDir := filepath.Clean(g.cfg.OutputDir)

	outputAbs, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolve output directory: %w", err)
	}

	if filepath.Dir(outputAbs) == outputAbs {
		return fmt.Errorf("output directory cannot be a filesystem root")
	}

	sourceDirs := []string{
		g.cfg.TemplatesDir,
		g.cfg.ContentDir,
		g.cfg.StaticDir,
	}

	for _, sourceDir := range sourceDirs {
		sourceAbs, err := filepath.Abs(filepath.Clean(sourceDir))
		if err != nil {
			return fmt.Errorf("resolve source directory %q: %w", sourceDir, err)
		}

		if pathsOverlap(outputAbs, sourceAbs) {
			return fmt.Errorf(
				"output directory %q overlaps source directory %q",
				outputDir,
				sourceDir,
			)
		}
	}

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("remove previous output: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	g.cfg.OutputDir = outputDir

	return nil
}

func pathsOverlap(first, second string) bool {
	return directoryContains(first, second) ||
		directoryContains(second, first)
}

func directoryContains(parent, child string) bool {
	relative, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	if relative == "." {
		return true
	}

	if relative == ".." || filepath.IsAbs(relative) {
		return false
	}

	return !strings.HasPrefix(
		relative,
		".."+string(os.PathSeparator),
	)
}

func (g *Generator) pageDefinitions() []pageDefinition {
	pagesDir := filepath.Join(g.cfg.ContentDir, "pages")

	return []pageDefinition{
		{
			Slug: "about",
			Page: PageData{
				Title:       "About",
				MetaTitle:   "About — Yash Sharma",
				Description: "About Yash Sharma.",
			},
			ContentFile: filepath.Join(pagesDir, "about.md"),
		},
		{
			Slug: "expertise",
			Page: PageData{
				Title:       "Expertise",
				MetaTitle:   "Expertise — Yash Sharma",
				Description: "Professional interests and areas of expertise.",
			},
			ContentFile: filepath.Join(pagesDir, "expertise.md"),
		},
		{
			Slug: "experience",
			Page: PageData{
				Title:       "Experience",
				MetaTitle:   "Experience — Yash Sharma",
				Description: "Selected work and professional experience.",
			},
			ContentFile: filepath.Join(pagesDir, "experience.md"),
		},
		{
			Slug: "contact",
			Page: PageData{
				Title:       "Contact",
				MetaTitle:   "Contact — Yash Sharma",
				Description: "Contact Yash Sharma.",
				Email:       "yash27062004@gmail.com",
			},
			ContentFile: filepath.Join(pagesDir, "contact.md"),
		},
	}
}

func (g *Generator) generateHome() error {
	page := PageData{
		Title:       "Home",
		MetaTitle:   "Yash Sharma — Software Engineer",
		Description: "Yash Sharma, Software Engineer.",
	}

	err := render.RenderPage(
		[]string{
			filepath.Join(g.cfg.TemplatesDir, "layout.html"),
			filepath.Join(g.cfg.TemplatesDir, "home.html"),
		},
		filepath.Join(g.cfg.OutputDir, "index.html"),
		page,
	)
	if err != nil {
		return fmt.Errorf("generate homepage: %w", err)
	}

	return nil
}

func (g *Generator) generatePages(pages []pageDefinition) error {
	for _, definition := range pages {
		body, err := content.ReadMarkdown(definition.ContentFile)
		if err != nil {
			return fmt.Errorf(
				"read %s page content: %w",
				definition.Slug,
				err,
			)
		}

		page := definition.Page
		page.Content = body

		destination := filepath.Join(
			g.cfg.OutputDir,
			definition.Slug,
			"index.html",
		)

		err = render.RenderPage(
			[]string{
				filepath.Join(g.cfg.TemplatesDir, "layout.html"),
				filepath.Join(g.cfg.TemplatesDir, "page.html"),
			},
			destination,
			page,
		)
		if err != nil {
			return fmt.Errorf(
				"generate %s page: %w",
				definition.Slug,
				err,
			)
		}
	}

	return nil
}

func (g *Generator) generateBlog(posts []content.Post) error {
	page := BlogData{
		PageData: PageData{
			Title:       "Blog",
			MetaTitle:   "Blog — Yash Sharma",
			Description: "Articles, observations, and ideas worth sharing.",
		},
		Posts: posts,
	}

	err := render.RenderPage(
		[]string{
			filepath.Join(g.cfg.TemplatesDir, "layout.html"),
			filepath.Join(g.cfg.TemplatesDir, "blog.html"),
		},
		filepath.Join(g.cfg.OutputDir, "blog", "index.html"),
		page,
	)
	if err != nil {
		return fmt.Errorf("generate blog listing: %w", err)
	}

	return nil
}

func (g *Generator) generateArticles(posts []content.Post) error {
	for _, post := range posts {
		page := ArticleData{
			PageData: PageData{
				Title:       post.Title,
				MetaTitle:   post.Title + " — Yash Sharma",
				Description: post.Summary,
			},
			Post: post,
		}

		destination := filepath.Join(
			g.cfg.OutputDir,
			"blog",
			post.Slug,
			"index.html",
		)

		err := render.RenderPage(
			[]string{
				filepath.Join(g.cfg.TemplatesDir, "layout.html"),
				filepath.Join(g.cfg.TemplatesDir, "post.html"),
			},
			destination,
			page,
		)
		if err != nil {
			return fmt.Errorf(
				"generate article %q: %w",
				post.Slug,
				err,
			)
		}
	}

	return nil
}

func (g *Generator) generateNotFound() error {
	path := filepath.Join(
		g.cfg.ContentDir,
		"pages",
		"not-found.md",
	)

	body, err := content.ReadMarkdown(path)
	if err != nil {
		return fmt.Errorf("read 404 content: %w", err)
	}

	page := PageData{
		Title:       "Page not found",
		MetaTitle:   "Page not found — Yash Sharma",
		Description: "The requested page could not be found.",
		Content:     body,
	}

	err = render.RenderPage(
		[]string{
			filepath.Join(g.cfg.TemplatesDir, "layout.html"),
			filepath.Join(g.cfg.TemplatesDir, "page.html"),
		},
		filepath.Join(g.cfg.OutputDir, "404.html"),
		page,
	)
	if err != nil {
		return fmt.Errorf("generate 404 page: %w", err)
	}

	return nil
}

func (g *Generator) sitemapPaths(
	pages []pageDefinition,
	posts []content.Post,
) []string {
	paths := []string{"/", "/blog/"}

	for _, page := range pages {
		paths = append(paths, "/"+page.Slug+"/")
	}

	for _, post := range posts {
		paths = append(paths, "/blog/"+post.Slug+"/")
	}

	return paths
}
