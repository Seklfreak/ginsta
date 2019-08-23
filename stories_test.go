package ginsta

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGinsta_StoriesByID(t *testing.T) {
	ginsta := NewGinsta(&http.Client{
		Timeout: 30 * time.Second,
	}, []string{"4240363941%3AhCFJFzEqfjcNFo%3A29"})

	posts, err := ginsta.StoriesByID(context.Background(), "1460871211")
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, post := range posts {
		fmt.Println(post.ID)
		fmt.Println(post.TakenAt)
	}
}
