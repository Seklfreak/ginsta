package ginsta

import (
	"context"
	"encoding/json"
	"errors"
)

func (g *Ginsta) UserByUsername(ctx context.Context, username string) (*User, error) {
	profile, err := g.userRawProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
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
					FullName        string `json:"full_name"`
					ID              string `json:"id"`
					IsPrivate       bool   `json:"is_private"`
					IsVerified      bool   `json:"is_verified"`
					ProfilePicURLHd string `json:"profile_pic_url_hd"`
					Username        string `json:"username"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
	RHXGIS string `json:"rhx_gis"`
}
