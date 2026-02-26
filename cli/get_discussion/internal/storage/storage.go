package storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/internal/discussion"
	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/pkg/urlutil"
)

var slugRe = regexp.MustCompile(`[^a-z0-9_\-]+`)

func slugifyTitle(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "_")
	slug = slugRe.ReplaceAllString(slug, "")
	slug = strings.Trim(slug, "_")
	if slug == "" {
		return "discussion"
	}
	return slug
}

func yamlEscape(value string) string {
	if value == "" {
		return ""
	}
	needsQuotes := strings.ContainsAny(value, ":#\n\r\t\"") ||
		strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ")
	if needsQuotes {
		escaped := strings.ReplaceAll(value, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		return fmt.Sprintf(`"%s"`, escaped)
	}
	return value
}

func buildFrontMatter(d *discussion.Discussion) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %s\n", yamlEscape(d.Title))
	fmt.Fprintf(&b, "link: %s\n", yamlEscape(d.Link))
	fmt.Fprintf(&b, "author: %s\n", yamlEscape(d.Author))
	fmt.Fprintf(&b, "comments: %s\n", yamlEscape(d.Comments))
	fmt.Fprintf(&b, "published_date: %s\n", yamlEscape(d.PublishedDate))
	b.WriteString("---\n\n")
	return b.String()
}

func readFrontMatter(path string) map[string]string {
	meta := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return meta
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return meta
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "---" {
			break
		}
		if line == "" || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		meta[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return meta
}

func LoadExistingLinks(outputDir string) map[string]string {
	links := map[string]string{}
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return links
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(outputDir, e.Name())
		meta := readFrontMatter(path)
		if link, ok := meta["link"]; ok && link != "" {
			links[urlutil.CanonicalizeURL(link)] = path
		}
	}
	return links
}

func ensureUniquePath(outputDir, baseName string, existingPaths map[string]struct{}) string {
	candidate := filepath.Join(outputDir, baseName+".md")
	if _, dup := existingPaths[candidate]; !dup {
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
	for i := 2; ; i++ {
		candidate = filepath.Join(outputDir, fmt.Sprintf("%s_%d.md", baseName, i))
		if _, dup := existingPaths[candidate]; !dup {
			if _, err := os.Stat(candidate); os.IsNotExist(err) {
				return candidate
			}
		}
	}
}

func SaveDiscussion(d *discussion.Discussion, outputDir string, existingByLink map[string]string) (string, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	linkKey := urlutil.CanonicalizeURL(d.Link)
	path, exists := existingByLink[linkKey]
	if !exists {
		slug := slugifyTitle(d.Title)
		existing := map[string]struct{}{}
		for _, v := range existingByLink {
			existing[v] = struct{}{}
		}
		path = ensureUniquePath(outputDir, slug, existing)
		existingByLink[linkKey] = path
	}

	content := buildFrontMatter(d) + strings.TrimSpace(d.ContentMD) + "\n"
	return path, os.WriteFile(path, []byte(content), 0o644)
}

func LoadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if key != "" {
			if _, already := os.LookupEnv(key); !already {
				os.Setenv(key, val)
			}
		}
	}
}
