// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "time"

// YouTubeVideo represents a YouTube video fetched from Invidious.
type YouTubeVideo struct {
	VideoID      string
	Title        string
	Author       string
	AuthorID     string
	AuthorAvatar string
	ViewCount    int
	LikeCount    int
	Duration     int // seconds, 0 for live streams
	PublishedAt  time.Time
	ThumbnailURL string
	IsLive       bool
}
