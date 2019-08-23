package ginsta

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGinsta_PostByShortcode(t *testing.T) {
	ginsta := NewGinsta(&http.Client{
		Timeout: 30 * time.Second,
	}, nil)

	post, err := ginsta.PostByShortcode(context.Background(), "Bv_XVAQBP8s")
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(post.ID)
	fmt.Println(post.Shortcode)
	fmt.Println(post.AuthorID)
	fmt.Println(post.AuthorUsername)
	for _, mediaURL := range post.MediaURLs {
		fmt.Println(mediaURL)
	}
}
