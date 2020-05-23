package ginsta

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (g *Ginsta) PostsByID(ctx context.Context, id string) ([]*Post, error) {
	endpoint, _ := endpointPosts(id)

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

	var posts postsData
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, err
	}

	if posts.Status != "ok" {
		return nil, errors.New("failure parsing posts")
	}

	var respondPosts []*Post

	if len(posts.Data.User.EdgeOwnerToTimelineMedia.Edges) == 0 {
		return respondPosts, nil
	}

	var doNotTryDetails bool
	for _, item := range posts.Data.User.EdgeOwnerToTimelineMedia.Edges {
		post := &Post{
			ID:        item.Node.ID,
			Shortcode: item.Node.Shortcode,
			TakenAt:   time.Unix(item.Node.TakenAtTimestamp, 0).UTC(),
			Comment:   item.Node.EdgeMediaToComment.Count,
			Likes:     item.Node.EdgeMediaPreviewLike.Count,
			MediaURLs: []string{item.Node.DisplayURL},
			Video:     item.Node.IsVideo,
		}

		if len(item.Node.EdgeMediaToCaption.Edges) > 0 {
			post.Caption = item.Node.EdgeMediaToCaption.Edges[0].Node.Text
		}

		if (item.Node.Typename == "GraphSidecar" ||
			item.Node.Typename == "GraphVideo") && !doNotTryDetails {
			postDetailed, err := g.PostByShortcode(ctx, item.Node.Shortcode)
			if err != nil {
				doNotTryDetails = true
				fmt.Printf("failure loading post %s, post will have less metadata\n", item.Node.Shortcode)
			} else {
				post = postDetailed
			}
		}

		respondPosts = append(respondPosts, post)
	}

	return respondPosts, nil
}

type postsData struct {
	Data struct {
		User struct {
			EdgeOwnerToTimelineMedia struct {
				Count int `json:"count"`
				Edges []struct {
					Node struct {
						ID                 string `json:"id"`
						Typename           string `json:"__typename"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int `json:"count"`
						} `json:"edge_media_to_comment"`
						CommentsDisabled bool  `json:"comments_disabled"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL           string `json:"display_url"`
						EdgeMediaPreviewLike struct {
							Count int `json:"count"`
						} `json:"edge_media_preview_like"`
						Owner struct {
							ID string `json:"id"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo bool `json:"is_video"`
					} `json:"node,omitempty"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}
