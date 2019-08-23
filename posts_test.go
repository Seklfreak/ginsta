package ginsta

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGinsta_PostsByID(t *testing.T) {
	ginsta := NewGinsta(&http.Client{
		Timeout: 30 * time.Second,
	}, nil)

	posts, err := ginsta.PostsByID(context.Background(), "2103418442")
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, post := range posts {
		fmt.Println(post.Shortcode)
		fmt.Println(post.TakenAt)
	}
}
