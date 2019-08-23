package ginsta

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGinsta_UserByUsername(t *testing.T) {
	ginsta := NewGinsta(&http.Client{
		Timeout: 30 * time.Second,
	}, nil)

	user, err := ginsta.UserByUsername(context.Background(), "taeri__taeri")
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(user.Username)
	for _, video := range user.Videos {
		fmt.Println(video.ID)
		fmt.Println(video.Product) // igtv
		fmt.Println(video.Shortcode)
		fmt.Println(video.Title)
	}
	fmt.Println(user.ID)
}
