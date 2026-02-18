package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WgerClient struct {
	baseURL string
	http    *http.Client
}

func NewWgerClient() *WgerClient {
	return &WgerClient{
		baseURL: "https://wger.de",
		http:    &http.Client{Timeout: 20 * time.Second},
	}
}

type wgerPaginatedResponse[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type WgerExerciseInfo struct {
	ID               int       `json:"id"`
	UUID             string    `json:"uuid"`
	Category         wgerNamed  `json:"category"`
	Muscles          []wgerNamed `json:"muscles"`
	MusclesSecondary []wgerNamed `json:"muscles_secondary"`
	Equipment        []wgerNamed `json:"equipment"`
	Images           []WgerExerciseImage `json:"images"`
	Translations     []WgerExerciseTranslation `json:"translations"`
	Videos           []WgerExerciseVideo `json:"videos"`
}

type wgerNamed struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WgerExerciseTranslation struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    int    `json:"language"`
}

type WgerExerciseImage struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	Image    string `json:"image"`
	IsMain   bool   `json:"is_main"`
	Exercise int    `json:"exercise"`
}

// WgerExerciseVideo is the nested exerciseinfo representation.
// The /api/v2/video/ endpoint uses the same field name (video).
type WgerExerciseVideo struct {
	ID           int    `json:"id"`
	UUID         string `json:"uuid"`
	Video        string `json:"video"`
	IsMain       bool   `json:"is_main"`
	Exercise     int    `json:"exercise"`
	ExerciseUUID string `json:"exercise_uuid"`
}

func (c *WgerClient) ListExerciseInfo(ctx context.Context, limit, offset int, languageCode string) (wgerPaginatedResponse[WgerExerciseInfo], error) {
	endpoint, err := url.Parse(c.baseURL + "/api/v2/exerciseinfo/")
	if err != nil {
		return wgerPaginatedResponse[WgerExerciseInfo]{}, err
	}

	q := endpoint.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	if languageCode != "" {
		q.Set("language__code", languageCode)
	}
	endpoint.RawQuery = q.Encode()

	var decoded wgerPaginatedResponse[WgerExerciseInfo]
	if err := c.doJSON(ctx, endpoint.String(), &decoded); err != nil {
		return wgerPaginatedResponse[WgerExerciseInfo]{}, err
	}
	return decoded, nil
}

func (c *WgerClient) ListVideos(ctx context.Context, limit, offset int) (wgerPaginatedResponse[WgerExerciseVideo], error) {
	endpoint, err := url.Parse(c.baseURL + "/api/v2/video/")
	if err != nil {
		return wgerPaginatedResponse[WgerExerciseVideo]{}, err
	}

	q := endpoint.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	endpoint.RawQuery = q.Encode()

	var decoded wgerPaginatedResponse[WgerExerciseVideo]
	if err := c.doJSON(ctx, endpoint.String(), &decoded); err != nil {
		return wgerPaginatedResponse[WgerExerciseVideo]{}, err
	}
	return decoded, nil
}

func (c *WgerClient) doJSON(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "S.P.A.R.T.A/1.0 (wger exercise library importer)")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("wger non-2xx: %d body=%q", resp.StatusCode, string(b))
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}
