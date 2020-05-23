package ginsta

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

func (g *Ginsta) PostByShortcode(ctx context.Context, shortcode string) (*Post, error) {
	endpoint, _ := endpointPost(shortcode)

	body, err := g.request(
		ctx,
		g.client,
		endpoint,
		"",
		false,
	)
	if err != nil {
		return nil, err
	}

	var post postSharedData
	err = json.Unmarshal([]byte(body), &post)
	if err != nil {
		return nil, err
	}

	if post.Status != "ok" {
		return nil, errors.New("failure parsing posts")
	}

	respPost := &Post{
		ID:             post.Data.ShortcodeMedia.ID,
		Shortcode:      post.Data.ShortcodeMedia.Shortcode,
		TakenAt:        time.Unix(post.Data.ShortcodeMedia.TakenAtTimestamp, 0).UTC(),
		MediaURLs:      getMediaURLs(&post),
		Comment:        post.Data.ShortcodeMedia.EdgeMediaToComment.Count,
		Likes:          post.Data.ShortcodeMedia.EdgeMediaPreviewLike.Count,
		Video:          post.Data.ShortcodeMedia.IsVideo,
		AuthorID:       post.Data.ShortcodeMedia.Owner.ID,
		AuthorUsername: post.Data.ShortcodeMedia.Owner.Username,
	}

	if len(post.Data.ShortcodeMedia.EdgeMediaToCaption.Edges) > 0 {
		respPost.Caption = post.Data.ShortcodeMedia.EdgeMediaToCaption.Edges[0].Node.Text
	}

	return respPost, nil
}

type Post struct {
	ID             string
	Shortcode      string
	Caption        string
	TakenAt        time.Time
	MediaURLs      []string
	Comment        int
	Likes          int
	Video          bool
	AuthorID       string
	AuthorUsername string
}

func getMediaURLs(post *postSharedData) []string {
	var urls []string

	if post.Data.ShortcodeMedia.VideoURL != "" {
		urls = append(urls,
			post.Data.ShortcodeMedia.VideoURL,
		)

		return urls
	}

	if len(post.Data.ShortcodeMedia.EdgeSidecarToChildren.Edges) > 0 {
		for _, item := range post.Data.ShortcodeMedia.EdgeSidecarToChildren.Edges {
			if item.Node.VideoURL != "" {
				urls = append(urls, item.Node.VideoURL)
				continue
			}

			urls = append(urls, getBestDisplayResource(item.Node.DisplayResources))
		}

		return urls
	}

	if len(post.Data.ShortcodeMedia.DisplayResources) > 0 {
		urls = append(urls, getBestDisplayResource(
			post.Data.ShortcodeMedia.DisplayResources,
		))

		return urls
	}

	return nil

}

type postSharedData struct {
	Data struct {
		ShortcodeMedia struct {
			ID                 string            `json:"id"`
			Shortcode          string            `json:"shortcode"`
			DisplayResources   []displayResource `json:"display_resources"`
			VideoURL           string            `json:"video_url"`
			IsVideo            bool              `json:"is_video"`
			EdgeMediaToCaption struct {
				Edges []struct {
					Node struct {
						Text string `json:"text"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_caption"`
			EdgeMediaToComment struct {
				Count int `json:"count"`
			} `json:"edge_media_to_comment"`
			TakenAtTimestamp     int64 `json:"taken_at_timestamp"`
			EdgeMediaPreviewLike struct {
				Count int `json:"count"`
			} `json:"edge_media_preview_like"`
			EdgeSidecarToChildren struct {
				Edges []struct {
					Node struct {
						Typename   string `json:"__typename"`
						ID         string `json:"id"`
						Shortcode  string `json:"shortcode"`
						Dimensions struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayResources []displayResource `json:"display_resources"`
						IsVideo          bool              `json:"is_video"`
						VideoURL         string            `json:"video_url"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_sidecar_to_children"`
			Owner struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			} `json:"owner"`
		} `json:"shortcode_media"`
	} `json:"data"`
	Status string `json:"status"`
}
