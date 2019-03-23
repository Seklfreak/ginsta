package ginsta

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

func (g *Ginsta) UserByUsername(ctx context.Context, username string) (*User, error) {
	profile, err := g.userRawProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	var videos []*Video
	for _, item := range profile.EntryData.ProfilePage[0].Graphql.User.EdgeFelixVideoTimeline.Edges {
		video := &Video{
			Product:    item.Node.ProductType,
			ID:         item.Node.ID,
			Shortcode:  item.Node.Shortcode,
			Title:      item.Node.Title,
			DisplayURL: item.Node.DisplayURL,
			TakentAt:   time.Unix(item.Node.TakenAtTimestamp, 0),
			Comments:   item.Node.EdgeMediaToComment.Count,
			Likes:      item.Node.EdgeLikedBy.Count,
			Published:  item.Node.IsPublished,
			Duration:   item.Node.VideoDuration,
		}

		if len(item.Node.EdgeMediaToCaption.Edges) > 0 {
			video.Caption = item.Node.EdgeMediaToCaption.Edges[0].Node.Text
		}

		videos = append(videos, video)
	}

	return &User{
		Username:      profile.EntryData.ProfilePage[0].Graphql.User.Username,
		FullName:      profile.EntryData.ProfilePage[0].Graphql.User.FullName,
		ID:            profile.EntryData.ProfilePage[0].Graphql.User.ID,
		IsPrivate:     profile.EntryData.ProfilePage[0].Graphql.User.IsPrivate,
		IsVerified:    profile.EntryData.ProfilePage[0].Graphql.User.IsVerified,
		ProfilePicURL: profile.EntryData.ProfilePage[0].Graphql.User.ProfilePicURLHd,
		Followers:     profile.EntryData.ProfilePage[0].Graphql.User.EdgeFollowedBy.Count,
		Followings:    profile.EntryData.ProfilePage[0].Graphql.User.EdgeFollow.Count,
		Videos:        videos,
	}, nil
}

func (g *Ginsta) userRawProfileByUsername(ctx context.Context, username string) (*profileSharedData, error) {
	body, err := g.request(
		ctx,
		g.client,
		endpointProfile(username),
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

	var profile profileSharedData
	err = json.Unmarshal([]byte(sharedData), &profile)
	if err != nil {
		return nil, err
	}

	if len(profile.EntryData.ProfilePage) == 0 ||
		profile.EntryData.ProfilePage[0].Graphql.User.ID == "" {
		return nil, errors.New("failure parsing profile")
	}

	return &profile, nil
}

type User struct {
	Username      string
	FullName      string
	ID            string
	IsPrivate     bool
	IsVerified    bool
	ProfilePicURL string
	Followers     int
	Followings    int
	Videos        []*Video
}

type Video struct {
	Product    string
	ID         string
	Shortcode  string
	Title      string
	Caption    string
	DisplayURL string
	TakentAt   time.Time
	Comments   int
	Likes      int
	Published  bool
	Duration   float64
}

type profileSharedData struct {
	EntryData struct {
		ProfilePage []struct {
			Graphql struct {
				User struct {
					EdgeFollowedBy struct {
						Count int `json:"count"`
					} `json:"edge_followed_by"`
					EdgeFollow struct {
						Count int `json:"count"`
					} `json:"edge_follow"`
					FullName               string `json:"full_name"`
					ID                     string `json:"id"`
					IsPrivate              bool   `json:"is_private"`
					IsVerified             bool   `json:"is_verified"`
					ProfilePicURLHd        string `json:"profile_pic_url_hd"`
					Username               string `json:"username"`
					EdgeFelixVideoTimeline struct {
						Count int `json:"count"`
						Edges []struct {
							Node struct {
								Typename           string `json:"__typename"`
								ID                 string `json:"id"`
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
								TakenAtTimestamp int64  `json:"taken_at_timestamp"`
								DisplayURL       string `json:"display_url"`
								EdgeLikedBy      struct {
									Count int `json:"count"`
								} `json:"edge_liked_by"`
								IsPublished   bool    `json:"is_published"`
								ProductType   string  `json:"product_type"`
								Title         string  `json:"title"`
								VideoDuration float64 `json:"video_duration"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_felix_video_timeline"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
	RHXGIS string `json:"rhx_gis"`
}
