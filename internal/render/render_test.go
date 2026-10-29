package render

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderPage(t *testing.T) {
	root := t.TempDir()

	layout := writeFixture(t, root, "layout.html", `{{define "layout"}}<!doctype html>
<html>
<head><title>{{.Title}}</title></head>
<body>{{template "content" .}}</body>
</html>{{end}}`)

	content := writeFixture(t, root, "content.html", `{{define "content"}}
<h1>{{.Title}}</h1>
<div>{{.Content}}</div>
{{end}}`)

	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   "Example Page",
		Content: template.HTML("<p>Rendered content.</p>"),
	}

	destination := filepath.Join(root, "dist", "index.html")

	if err := RenderPage(
		[]string{layout, content},
		destination,
		data,
	); err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	output := readFixture(t, destination)

	for _, expected := range []string{
		"<title>Example Page</title>",
		"<h1>Example Page</h1>",
		"<p>Rendered content.</p>",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("output does not contain %q", expected)
		}
	}
}

func TestRenderPageEscapesOrdinaryStringValues(t *testing.T) {
	root := t.TempDir()

	layout := writeFixture(t, root, "layout.html", `{{define "layout"}}{{.Title}}{{end}}`)

	content := writeFixture(t, root, "content.html", `{{define "content"}}Unused{{end}}`)

	destination := filepath.Join(root, "dist", "index.html")

	err := RenderPage(
		[]string{layout, content},
		destination,
		struct {
			Title string
		}{
			Title: "<script>alert(1)</script>",
		},
	)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	output := readFixture(t, destination)

	if strings.Contains(output, "<script>") {
		t.Fatalf("ordinary string was not escaped: %s", output)
	}

	if !strings.Contains(output, "&lt;script&gt;") {
		t.Errorf("escaped text missing from output: %s", output)
	}
}

func TestRenderPageRejectsEmptyTemplateList(t *testing.T) {
	err := RenderPage(nil, filepath.Join(t.TempDir(), "index.html"), nil)

	if err == nil {
		t.Fatal("RenderPage() expected an error")
	}
}

func TestRenderPageReportsMissingTemplate(t *testing.T) {
	root := t.TempDir()

	err := RenderPage(
		[]string{filepath.Join(root, "missing.html")},
		filepath.Join(root, "dist", "index.html"),
		nil,
	)

	if err == nil {
		t.Fatal("RenderPage() expected an error")
	}

	if !strings.Contains(err.Error(), "parse templates") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRenderPageReportsTemplateExecutionError(t *testing.T) {
	root := t.TempDir()

	layout := writeFixture(
		t,
		root,
		"layout.html",
		`{{define "layout"}}{{.Missing.Field}}{{end}}`,
	)

	content := writeFixture(
		t,
		root,
		"content.html",
		`{{define "content"}}Content{{end}}`,
	)

	err := RenderPage(
		[]string{layout, content},
		filepath.Join(root, "dist", "index.html"),
		struct{ Title string }{Title: "Example"},
	)

	if err == nil {
		t.Fatal("RenderPage() expected an execution error")
	}

	if !strings.Contains(err.Error(), "render template") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCopyDirCopiesNestedFiles(t *testing.T) {
	root := t.TempDir()

	source := filepath.Join(root, "static")
	destination := filepath.Join(root, "dist", "static")

	writeFixture(
		t,
		source,
		filepath.Join("css", "style.css"),
		"body { color: black; }\n",
	)

	if err := CopyDir(source, destination); err != nil {
		t.Fatalf("CopyDir() error = %v", err)
	}

	copiedFile := filepath.Join(destination, "css", "style.css")
	assertFileContent(t, copiedFile, "body { color: black; }\n")
}

func TestCopyDirRejectsMissingSource(t *testing.T) {
	root := t.TempDir()

	err := CopyDir(
		filepath.Join(root, "missing"),
		filepath.Join(root, "dist"),
	)

	if err == nil {
		t.Fatal("CopyDir() expected an error")
	}
}

func TestCopyDirRejectsOverlappingDirectories(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "static")

	if err := os.MkdirAll(source, 0755); err != nil {
		t.Fatal(err)
	}

	err := CopyDir(source, filepath.Join(source, "backup"))
	if err == nil {
		t.Fatal("CopyDir() expected an error for overlapping paths")
	}
}

func writeFixture(
	t *testing.T,
	root string,
	name string,
	source string,
) string {
	t.Helper()

	path := filepath.Join(root, name)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	return path
}

func readFixture(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	return string(data)
}

func assertFileContent(t *testing.T, path, expected string) {
	t.Helper()

	actual := readFixture(t, path)

	if actual != expected {
		t.Errorf("file content = %q, want %q", actual, expected)
	}
}
