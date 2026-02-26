package discussion

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/internal/client"
	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/pkg/urlutil"
)

// ExtractDiscussionLinksFromHTML parses anchor hrefs from raw HTML bytes.
func ExtractDiscussionLinksFromHTML(body []byte, base string) []string {
	seen := map[string]struct{}{}
	var out []string

	hrefRegex := regexp.MustCompile(`(?i)href=["']([^"']+)["']`)
	for _, m := range hrefRegex.FindAllSubmatch(body, -1) {
		href := string(m[1])
		if !strings.Contains(href, "/discussions/") && !strings.Contains(href, "/discussion/") {
			continue
		}
		baseURL, _ := url.Parse(base)
		ref, err := url.Parse(href)
		if err != nil {
			continue
		}
		abs := urlutil.CanonicalizeURL(baseURL.ResolveReference(ref).String())
		if _, ok := seen[abs]; !ok {
			seen[abs] = struct{}{}
			out = append(out, abs)
		}
	}
	return out
}

// extractTitleFromHTML returns a best-effort title from raw HTML.
func extractTitleFromHTML(body []byte) string {
	h1Re := regexp.MustCompile(`(?is)<h1[^>]*>(.*?)</h1>`)
	if m := h1Re.FindSubmatch(body); m != nil {
		return strings.TrimSpace(stripTags(string(m[1])))
	}
	ogRe := regexp.MustCompile(`(?i)<meta[^>]+property=["']og:title["'][^>]+content=["']([^"']+)["']`)
	if m := ogRe.FindSubmatch(body); m != nil {
		return strings.TrimSpace(string(m[1]))
	}
	return "untitled_discussion"
}

// stripTags removes HTML tags from a string.
func stripTags(s string) string {
	return regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, "")
}

// htmlToMarkdown does a simplistic conversion: strips tags, collapses whitespace.
func htmlToMarkdown(body []byte) string {
	s := string(body)
	// Strip scripts/styles.
	s = regexp.MustCompile(`(?is)<(script|style|noscript)[^>]*>.*?</(script|style|noscript)>`).ReplaceAllString(s, "")
	// Headings.
	for i := 6; i >= 1; i-- {
		tag := fmt.Sprintf("h%d", i)
		prefix := strings.Repeat("#", i) + " "
		s = regexp.MustCompile(fmt.Sprintf(`(?is)<%s[^>]*>(.*?)</%s>`, tag, tag)).
			ReplaceAllStringFunc(s, func(m string) string {
				inner := regexp.MustCompile(fmt.Sprintf(`(?is)<%s[^>]*>(.*?)</%s>`, tag, tag)).FindStringSubmatch(m)
				if inner != nil {
					return "\n" + prefix + stripTags(inner[1]) + "\n"
				}
				return m
			})
	}
	// Paragraphs and line breaks.
	s = regexp.MustCompile(`(?i)<br\s*/?>|</p>|</div>|</li>`).ReplaceAllString(s, "\n")
	// Links.
	s = regexp.MustCompile(`(?is)<a[^>]+href=["']([^"']+)["'][^>]*>(.*?)</a>`).
		ReplaceAllString(s, "[$2]($1)")
	// Bold/italic.
	s = regexp.MustCompile(`(?is)<(strong|b)>(.*?)</(strong|b)>`).ReplaceAllString(s, "**$2**")
	s = regexp.MustCompile(`(?is)<(em|i)>(.*?)</(em|i)>`).ReplaceAllString(s, "_$2_")
	// Strip remaining tags.
	s = stripTags(s)
	// Collapse excessive blank lines.
	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

func BuildDiscussionFromHTML(c *client.Client, rawURL string) (*Discussion, error) {
	body, err := c.FetchBody(rawURL, nil)
	if err != nil {
		return nil, err
	}
	title := extractTitleFromHTML(body)
	contentMD := htmlToMarkdown(body)
	return &Discussion{
		Title:     title,
		Link:      urlutil.CanonicalizeURL(rawURL),
		ContentMD: contentMD,
	}, nil
}
