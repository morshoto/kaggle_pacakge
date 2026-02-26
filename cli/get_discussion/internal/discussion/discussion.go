package discussion

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/internal/api"
	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/internal/client"
	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/pkg/urlutil"
)

// Discussion holds all metadata and content for a single Kaggle discussion.
type Discussion struct {
	Title         string
	Link          string
	Author        string
	Comments      string
	PublishedDate string
	ContentMD     string
}

func BuildDiscussionFromAPI(c *client.Client, rawURL string, topicID int) (*Discussion, error) {
	// Warm up cookies.
	_, _ = c.FetchBody(rawURL, nil)

	topicResp, err := api.FetchTopicData(c, topicID)
	if err != nil {
		return nil, err
	}
	t := topicResp.ForumTopic
	if t.Name == "" {
		return nil, fmt.Errorf("empty forumTopic for topic_id=%d", topicID)
	}

	msgResp, err := api.FetchTopicMessages(c, topicID)
	if err != nil {
		return nil, err
	}

	// Find the first / main message.
	var contentMD string
	for _, m := range msgResp.Comments {
		if m.ID == t.FirstMessageID || contentMD == "" {
			contentMD = m.RawMarkdown
			if contentMD == "" {
				contentMD = m.Content
			}
			if m.ID == t.FirstMessageID {
				break
			}
		}
	}

	link := t.URL
	if link == "" {
		link = rawURL
	}
	if !strings.HasPrefix(link, "http") {
		base, _ := url.Parse("https://www.kaggle.com")
		ref, _ := url.Parse(link)
		link = base.ResolveReference(ref).String()
	}
	link = urlutil.CanonicalizeURL(link)

	author := t.AuthorUserDisplayName
	if author == "" {
		author = t.AuthorUserName
	}

	comments := ""
	if t.TotalMessages != nil {
		comments = fmt.Sprint(*t.TotalMessages)
	}

	return &Discussion{
		Title:         urlutil.FirstNonEmpty(t.Name, "untitled_discussion"),
		Link:          link,
		Author:        author,
		Comments:      comments,
		PublishedDate: t.PostDate,
		ContentMD:     strings.TrimSpace(contentMD),
	}, nil
}

// IterDiscussions yields Discussion values for each URL, with API -> HTML fallback.
func IterDiscussions(urls []string, c *client.Client, delay time.Duration) <-chan *Discussion {
	ch := make(chan *Discussion)
	go func() {
		defer close(ch)
		for _, rawURL := range urls {
			topicID, hasID := urlutil.ExtractTopicID(rawURL)
			var d *Discussion
			var err error

			if hasID {
				d, err = BuildDiscussionFromAPI(c, rawURL, topicID)
				if err != nil {
					log.Printf("[warn] API failed for %s: %v — falling back to HTML", rawURL, err)
					d, err = BuildDiscussionFromHTML(c, rawURL)
				}
			} else {
				log.Printf("[warn] No topic ID detected in URL %s — using HTML parser", rawURL)
				d, err = BuildDiscussionFromHTML(c, rawURL)
			}

			if err != nil {
				log.Printf("[warn] Skipping %s: %v", rawURL, err)
				continue
			}
			if strings.TrimSpace(d.ContentMD) == "" {
				log.Printf("[warn] Empty content for %s", rawURL)
			}
			ch <- d
			if delay > 0 {
				time.Sleep(delay)
			}
		}
	}()
	return ch
}
