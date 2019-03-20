package ginsta

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	userAgent = "Instagram 10.26.0 (iPhone7,2; iOS 10_1_1; en_US; en-US; scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+"
)

var (
	endpointProfile = func(username string) string {
		return fmt.Sprintf("https://www.instagram.com/%s/", username)
	}

	endpointPost = func(shortcode string) string {
		return fmt.Sprintf("https://www.instagram.com/p/%s/", shortcode)
	}

	endpointPosts = func(id string) (string, string) {
		data := struct {
			ID    string `json:"id"`
			First int    `json:"first"`
		}{
			ID:    id,
			First: 15,
		}

		marshalled, _ := json.Marshal(data)

		return "https://www.instagram.com/graphql/query/?query_id=17888483320059182&variables=" +
			url.QueryEscape(string(marshalled)), string(marshalled)
	}

	endpointStories = func(id string) (string, string) {
		data := struct {
			ReelIDs            []string `json:"reel_ids"`
			PrecomposedOverlay bool     `json:"precomposed_overlay"`
		}{
			ReelIDs:            []string{id},
			PrecomposedOverlay: false,
		}

		marshalled, _ := json.Marshal(data)

		return "https://www.instagram.com/graphql/query/?query_id=17873473675158481&variables=" +
			url.QueryEscape(string(marshalled)), string(marshalled)
	}
)
