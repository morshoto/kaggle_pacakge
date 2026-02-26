package api

import (
	"fmt"
	"net/url"

	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/internal/client"
	"github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion/pkg/urlutil"
)

const (
	apiTopicURL       = "https://www.kaggle.com/api/i/discussions.DiscussionsService/GetForumTopicById"
	apiMessagesURL    = "https://www.kaggle.com/api/i/discussions.DiscussionsService/GetForumMessagesInTopic"
	apiCompetitionURL = "https://www.kaggle.com/api/i/competitions.CompetitionService/GetCompetition"
	apiTopicListURL   = "https://www.kaggle.com/api/i/discussions.DiscussionsService/GetTopicListByForumId"
)

func FetchTopicData(c *client.Client, topicID int) (*TopicResponse, error) {
	params := url.Values{"forumTopicId": {fmt.Sprint(topicID)}}
	var resp TopicResponse
	if err := c.FetchJSON(apiTopicURL, params, &resp); err != nil {
		return nil, err
	}
	c.LogInfo("Topic API ok topic_id=%d", topicID)
	return &resp, nil
}

func FetchTopicMessages(c *client.Client, topicID int) (*MessagesResponse, error) {
	var resp MessagesResponse
	if err := c.PostJSONDecode(apiMessagesURL, map[string]any{
		"topicId":                 topicID,
		"includeFirstForumMessage": true,
	}, &resp); err != nil {
		return nil, err
	}
	c.LogInfo("Messages API ok topic_id=%d count=%d", topicID, len(resp.Comments))
	return &resp, nil
}

func FetchCompetitionForumID(c *client.Client, competition string) (int, error) {
	params := url.Values{"competitionName": {competition}}
	var resp CompetitionResponse
	if err := c.FetchJSON(apiCompetitionURL, params, &resp); err != nil {
		return 0, err
	}
	if resp.ForumID == nil {
		return 0, fmt.Errorf("forumId missing for competition=%s", competition)
	}
	c.LogInfo("Competition API ok competition=%s forumId=%d", competition, *resp.ForumID)
	return *resp.ForumID, nil
}

func FetchTopicListByForumID(c *client.Client, forumID int, sortKey, timeKey string, limit int) ([]string, error) {
	var allURLs []string
	total := -1

	for page := 1; ; page++ {
		params := url.Values{
			"forumId": {fmt.Sprint(forumID)},
			"page":    {fmt.Sprint(page)},
		}
		if v, ok := urlutil.SortParam(sortKey); ok {
			params.Set("sort", v)
		}
		if v, ok := urlutil.TimeFilterParam(timeKey); ok {
			params.Set("time", v)
		}

		var resp TopicListResponse
		if err := c.FetchJSON(apiTopicListURL, params, &resp); err != nil {
			return allURLs, err
		}
		c.LogInfo("Topic list API ok forum_id=%d page=%d count=%d", forumID, page, resp.Count)

		if total < 0 {
			total = resp.Count
		}
		if len(resp.Topics) == 0 {
			break
		}

		for _, t := range resp.Topics {
			raw := t.TopicURL
			if raw == "" {
				raw = t.URL
			}
			if raw == "" {
				continue
			}
			abs, _ := url.Parse("https://www.kaggle.com")
			ref, err := url.Parse(raw)
			if err == nil {
				allURLs = append(allURLs, urlutil.CanonicalizeURL(abs.ResolveReference(ref).String()))
			}
		}

		if limit > 0 && len(allURLs) >= limit {
			allURLs = allURLs[:limit]
			break
		}

		if total >= 0 && len(allURLs) >= total {
			break
		}
	}
	return allURLs, nil
}
