package ginsta

import (
	"context"
	"encoding/json"
	"errors"
)

func (g *Ginsta) PostByShortcode(ctx context.Context, shortcode string) (*Post, error) {
	body, err := g.request(
		ctx,
		g.client,
		endpointPost(shortcode),
		"",
		false,
	)
	if err != nil {
		return nil, err
	}

	sharedData, err := extractSharedData(string(body))
	if err != nil {
		return nil, err
	}

	var post postSharedData
	err = json.Unmarshal([]byte(sharedData), &post)
	if err != nil {
		return nil, err
	}

	if len(post.EntryData.PostPage) == 0 {
		return nil, errors.New("failure parsing post")
	}

	respPost := &Post{
		ID:        post.EntryData.PostPage[0].Graphql.ShortcodeMedia.ID,
		Shortcode: post.EntryData.PostPage[0].Graphql.ShortcodeMedia.Shortcode,
		TakenAt:   post.EntryData.PostPage[0].Graphql.ShortcodeMedia.TakenAtTimestamp,
		MediaURLs: getMediaURLs(&post),
		Comment:   post.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeMediaToComment.Count,
		Likes:     post.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeMediaPreviewLike.Count,
		Video:     post.EntryData.PostPage[0].Graphql.ShortcodeMedia.IsVideo,
	}

	if len(post.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeMediaToCaption.Edges) > 0 {
		respPost.Caption = post.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeMediaToCaption.Edges[0].Node.Text
	}

	return respPost, nil
}

type Post struct {
	ID        string
	Shortcode string
	Caption   string
	TakenAt   int
	MediaURLs []string
	Comment   int
	Likes     int
	Video     bool
}

func getMediaURLs(post *postSharedData) []string {
	var urls []string
	if len(post.EntryData.PostPage) == 0 {
		return urls
	}

	if post.EntryData.PostPage[0].Graphql.ShortcodeMedia.VideoURL != "" {
		urls = append(urls,
			post.EntryData.PostPage[0].Graphql.ShortcodeMedia.VideoURL,
		)

		return urls
	}

	if len(post.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges) > 0 {
		for _, item := range post.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges {
			urls = append(urls, getBestDisplayResource(item.Node.DisplayResources))
		}

		return urls
	}

	if len(post.EntryData.PostPage[0].Graphql.ShortcodeMedia.DisplayResources) > 0 {
		urls = append(urls, getBestDisplayResource(
			post.EntryData.PostPage[0].Graphql.ShortcodeMedia.DisplayResources,
		))

		return urls
	}

	return nil

}

type postSharedData struct {
	EntryData struct {
		PostPage []struct {
			Graphql struct {
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
					TakenAtTimestamp     int `json:"taken_at_timestamp"`
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
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_sidecar_to_children"`
				} `json:"shortcode_media"`
			} `json:"graphql"`
		} `json:"PostPage"`
	} `json:"entry_data"`
}