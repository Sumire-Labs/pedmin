package rss

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mmcdole/gofeed"
	"github.com/s12kuma01/pedmin/store"
)

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

	for _, item := range newItems {
		msg := BuildFeedAnnouncement(feed.Title, item)
		if _, err := (*r.client).Rest.CreateMessage(feed.ChannelID, msg); err != nil {
			r.logger.Warn("failed to post feed item",
				slog.Int64("feed_id", feed.ID),
				slog.Any("error", err),
			)
		}
	}

	if err := r.store.MarkItemsSeen(feed.ID, newHashes); err != nil {
		return fmt.Errorf("failed to mark items seen: %w", err)
	}

	r.logger.Info("posted new feed items",
		slog.Int64("feed_id", feed.ID),
		slog.Int("count", len(newItems)),
	)
	return nil
}
