package ginsta

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

func (g *Ginsta) StoriesByID(ctx context.Context, id string) ([]*Story, error) {
	endpoint, data := endpointStories(id)

	gis, err := g.gis(ctx, data)
	if err != nil {
		return nil, err
	}

	body, err := g.request(
		ctx,
		g.client,
		endpoint,
		gis,
		true,
	)
	if err != nil {
		return nil, err
	}

	var stories storiesData
	err = json.Unmarshal(body, &stories)
	if err != nil {
		return nil, err
	}

	if len(stories.Data.ReelsMedia) == 0 {
		return nil, errors.New("failure parsing stories")
	}

	var respondStories []*Story
	for _, item := range stories.Data.ReelsMedia[0].Items {
		story := &Story{
			ID:      item.ID,
			TakenAt: time.Unix(item.TakenAtTimestamp, 0).UTC(),
			Video:   false,
		}

		if item.IsVideo {
			story.Video = true

			story.MediaURLs = append(story.MediaURLs, getBestDisplayResource(item.VideoResources))
		} else {

			story.MediaURLs = append(story.MediaURLs, getBestDisplayResource(item.DisplayResources))
		}

		if len(story.MediaURLs) == 0 {
			continue
		}

		respondStories = append(respondStories, story)
	}

	if len(respondStories) == 0 {
		return nil, errors.New("failure parsing story items")
	}

	return respondStories, nil
}

type Story struct {
	ID        string
	MediaURLs []string
	Video     bool
	TakenAt   time.Time
}

type storiesData struct {
	Data struct {
		ReelsMedia []struct {
			ID    string `json:"id"`
			Items []struct {
				Typename   string `json:"__typename"`
				ID         string `json:"id"`
				Dimensions struct {
					Height int `json:"height"`
					Width  int `json:"width"`
				} `json:"dimensions"`
				DisplayResources []displayResource `json:"display_resources"`
				IsVideo          bool              `json:"is_video"`
				VideoDuration    float64           `json:"video_duration,omitempty"`
				VideoResources   []displayResource `json:"video_resources,omitempty"`
				TakenAtTimestamp int64             `json:"taken_at_timestamp"`
			} `json:"items"`
		} `json:"reels_media"`
	} `json:"data"`
	Status string `json:"status"`
}
