package rss

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mmcdole/gofeed"
	"github.com/s12kuma01/pedmin/store"
)

const MaxFeedsPerGuild = 10

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func (r *RSS) AddFeed(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, url string) (*store.RSSFeed, error) {
	count, err := r.store.CountRSSFeeds(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to count feeds: %w", err)
	}
	if count >= MaxFeedsPerGuild {
		return nil, fmt.Errorf("フィード数が上限（%d件）に達しています", MaxFeedsPerGuild)
	}

	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("フィードの取得に失敗しました: %w", err)
	}

	feed := &store.RSSFeed{
		GuildID:   guildID,
		URL:       url,
		ChannelID: channelID,
		Title:     parsed.Title,
	}

	if err := r.store.CreateRSSFeed(feed); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, fmt.Errorf("このフィードは既に登録されています")
		}
		return nil, fmt.Errorf("failed to create feed: %w", err)
	}

	// Mark all existing items as seen to prevent flood
	var hashes []string
	for _, item := range parsed.Items {
		hashes = append(hashes, itemHash(item))
	}
	if len(hashes) > 0 {
		if err := r.store.MarkItemsSeen(feed.ID, hashes); err != nil {
			r.logger.Warn("failed to mark existing items as seen", slog.Any("error", err))
		}
	}

	// Post preview of the latest item
	if len(parsed.Items) > 0 {
		msg := BuildFeedAnnouncement(feed.Title, parsed.Items[0])
		if _, err := (*r.client).Rest.CreateMessage(channelID, msg); err != nil {
			r.logger.Warn("failed to post preview", slog.Any("error", err))
		}
	}

	return feed, nil
}

func (r *RSS) RemoveFeed(feedID int64, guildID snowflake.ID) error {
	return r.store.DeleteRSSFeed(feedID, guildID)
}

func (r *RSS) pollSingleFeed(ctx context.Context, feed store.RSSFeed) error {
	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(feed.URL, ctx)
	if err != nil {
		return fmt.Errorf("failed to parse feed: %w", err)
	}

	// Find new items (iterate in reverse to post oldest first)
	var newItems []*gofeed.Item
	var newHashes []string
	for i := len(parsed.Items) - 1; i >= 0; i-- {
		item := parsed.Items[i]
		hash := itemHash(item)
		seen, err := r.store.IsItemSeen(feed.ID, hash)
		if err != nil {
			return fmt.Errorf("failed to check seen: %w", err)
		}
		if !seen {
			newItems = append(newItems, item)
			newHashes = append(newHashes, hash)
		}
	}

	if len(newItems) == 0 {
		return nil
	}

	// Post new items
	for _, item := range newItems {
		msg := BuildFeedAnnouncement(feed.Title, item)
		if _, err := (*r.client).Rest.CreateMessage(feed.ChannelID, msg); err != nil {
			r.logger.Warn("failed to post feed item",
				slog.Int64("feed_id", feed.ID),
				slog.Any("error", err),
			)
		}
	}

	// Mark all new items as seen
	if err := r.store.MarkItemsSeen(feed.ID, newHashes); err != nil {
		return fmt.Errorf("failed to mark items seen: %w", err)
	}

	r.logger.Info("posted new feed items",
		slog.Int64("feed_id", feed.ID),
		slog.Int("count", len(newItems)),
	)
	return nil
}

func itemHash(item *gofeed.Item) string {
	key := item.GUID
	if key == "" {
		key = item.Link
	}
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h)
}

func stripHTML(s string) string {
	return strings.TrimSpace(htmlTagRe.ReplaceAllString(s, ""))
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func ephemeralV2(components ...discord.LayoutComponent) discord.MessageCreate {
	return discord.NewMessageCreateV2(components...).WithEphemeral(true)
}
