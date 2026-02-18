package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WikipediaClient struct {
	http *http.Client
}

func NewWikipediaClient() *WikipediaClient {
	return &WikipediaClient{
		http: &http.Client{Timeout: 12 * time.Second},
	}
}

type wikipediaSummaryResponse struct {
	Thumbnail *struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

func (c *WikipediaClient) GetPageThumbnailURL(ctx context.Context, title string) (string, error) {
	if strings.TrimSpace(title) == "" {
		return "", errors.New("empty title")
	}

	encoded := url.PathEscape(title)
	endpoint := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", encoded)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	// Wikipedia requires a descriptive UA.
	req.Header.Set("User-Agent", "S.P.A.R.T.A/1.0 (exercise library seeder)")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("wikipedia non-2xx: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var decoded wikipediaSummaryResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", err
	}
	if decoded.Thumbnail == nil || strings.TrimSpace(decoded.Thumbnail.Source) == "" {
		return "", nil
	}
	return decoded.Thumbnail.Source, nil
}
