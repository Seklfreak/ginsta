package ginsta

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	userAgent = "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)"
)

var (
	endpointProfile = func(username string) string {
		return fmt.Sprintf("https://www.instagram.com/%s/", username)
	}

	endpointPost = func(shortcode string) (string, string) {
		data := struct {
			Shortcode string `json:"shortcode"`
		}{
			Shortcode: shortcode,
		}

		marshalled, _ := json.Marshal(data)

		return "https://www.instagram.com/graphql/query/?query_hash=55a3c4bad29e4e20c20ff4cdfd80f5b4&variables=" +
			url.QueryEscape(string(marshalled)), string(marshalled)
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
