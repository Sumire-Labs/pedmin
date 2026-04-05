// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Sumire-Labs/pedmin/internal/model"
)

// InvidiousClient fetches YouTube video data from an Invidious instance.
type InvidiousClient struct {
	baseURL string
	http    *http.Client
}

// NewInvidiousClient creates a new InvidiousClient.
func NewInvidiousClient(baseURL string, timeout time.Duration) *InvidiousClient {
	return &InvidiousClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: timeout},
	}
}

type invidiousVideo struct {
	Title            string               `json:"title"`
	VideoID          string               `json:"videoId"`
	Author           string               `json:"author"`
	AuthorID         string               `json:"authorId"`
	AuthorThumbnails []invidiousThumbnail `json:"authorThumbnails"`
	ViewCount        int                  `json:"viewCount"`
	LikeCount        int                  `json:"likeCount"`
	LengthSeconds    int                  `json:"lengthSeconds"`
	Published        int64                `json:"published"`
	VideoThumbnails  []invidiousThumbnail `json:"videoThumbnails"`
	LiveNow          bool                 `json:"liveNow"`
}

type invidiousThumbnail struct {
	URL     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Quality string `json:"quality"`
}

// GetVideo fetches video data from the Invidious API.
func (c *InvidiousClient) GetVideo(ctx context.Context, videoID string) (*model.YouTubeVideo, error) {
	endpoint := fmt.Sprintf("%s/api/v1/videos/%s?fields=title,videoId,author,authorId,authorThumbnails,viewCount,likeCount,lengthSeconds,published,videoThumbnails,liveNow",
		c.baseURL, videoID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invidious API returned %d: %s", resp.StatusCode, string(body))
	}

	var iv invidiousVideo
	if err := json.Unmarshal(body, &iv); err != nil {
		return nil, fmt.Errorf("invalid invidious response: %w", err)
	}

	return &model.YouTubeVideo{
		VideoID:      iv.VideoID,
		Title:        iv.Title,
		Author:       iv.Author,
		AuthorID:     iv.AuthorID,
		AuthorAvatar: pickThumbnail(iv.AuthorThumbnails),
		ViewCount:    iv.ViewCount,
		LikeCount:    iv.LikeCount,
		Duration:     iv.LengthSeconds,
		PublishedAt:  time.Unix(iv.Published, 0),
		ThumbnailURL: pickVideoThumbnail(iv.VideoThumbnails),
		IsLive:       iv.LiveNow,
	}, nil
}

// pickThumbnail selects the best author thumbnail (prefer ~176px).
func pickThumbnail(thumbs []invidiousThumbnail) string {
	var best string
	for _, t := range thumbs {
		best = t.URL
		if t.Width >= 176 {
			return t.URL
		}
	}
	return best
}

// pickVideoThumbnail selects the best quality video thumbnail.
func pickVideoThumbnail(thumbs []invidiousThumbnail) string {
	preferred := []string{"maxresdefault", "sddefault", "high", "medium"}
	for _, q := range preferred {
		for _, t := range thumbs {
			if t.Quality == q {
				return t.URL
			}
		}
	}
	if len(thumbs) > 0 {
		return thumbs[0].URL
	}
	return ""
}
