package ginsta

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (g *Ginsta) request(
	ctx context.Context,
	client *http.Client,
	endpoint string,
	gis string,
	authenticated bool,
) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("user-agent", userAgent)
	req.Header.Set("x-ig-capabilities", "36oD")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.8")
	req.Header.Set("x-instagram-gis", gis)

	if authenticated {
		if len(g.sessionIDs) == 0 {
			return nil, errors.New("no session IDs passed")
		}

		req.AddCookie(&http.Cookie{
			Name:  "sessionid",
			Value: g.sessionIDs[g.random.Intn(len(g.sessionIDs))],
		})
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
