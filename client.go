package ginsta

import (
	"math/rand"
	"net/http"
	"time"
)

type Ginsta struct {
	client     *http.Client
	sessionIDs []string

	rhxgis string
	random *rand.Rand
}

func NewGinsta(client *http.Client, sessionIDs []string) *Ginsta {
	return &Ginsta{
		client:     client,
		sessionIDs: sessionIDs,

		random: rand.New(rand.NewSource(time.Now().Unix())),
	}
}
