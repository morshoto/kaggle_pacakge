package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/internal/discussion"
)

func TestSlugifyTitle(t *testing.T) {
	got := slugifyTitle("Hello, World!!")
	if got != "hello_world" {
		t.Fatalf("unexpected slug: %s", got)
	}
}

func TestFrontMatterRoundTrip(t *testing.T) {
	d := &discussion.Discussion{
		Title:         "My Title",
		Link:          "https://example.com",
		Author:        "Author",
		Comments:      "5",
		PublishedDate: "2024-01-01",
		ContentMD:     "Body",
	}
	content := buildFrontMatter(d) + "Body\n"

	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	meta := readFrontMatter(path)
	if meta["title"] != d.Title {
		t.Fatalf("title mismatch: %s", meta["title"])
	}
	if meta["link"] != d.Link {
		t.Fatalf("link mismatch: %s", meta["link"])
	}
}

func TestEnsureUniquePath(t *testing.T) {
	dir := t.TempDir()
	base := "discussion"
	first := filepath.Join(dir, base+".md")
	if err := os.WriteFile(first, []byte("x"), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	paths := map[string]struct{}{first: {}}
	got := ensureUniquePath(dir, base, paths)
	if got == first {
		t.Fatalf("expected unique path, got %s", got)
	}
}
