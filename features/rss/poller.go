package rss

import (
	"context"
	"log/slog"
	"time"
)

const PollInterval = 5 * time.Minute

func (r *RSS) StartPoller(ctx context.Context) {
	ctx, r.cancel = context.WithCancel(ctx)
	go r.pollLoop(ctx)
}

func (r *RSS) StopPoller() {
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *RSS) pollLoop(ctx context.Context) {
	r.logger.Info("rss poller started", slog.Duration("interval", PollInterval))

	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	// Prune old seen items on startup
	r.pruneOldItems()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("rss poller stopped")
			return
		case <-ticker.C:
			r.pollAllFeeds(ctx)
		}
	}
}

func (r *RSS) pollAllFeeds(ctx context.Context) {
	feeds, err := r.store.GetAllRSSFeeds()
	if err != nil {
		r.logger.Error("failed to get all rss feeds", slog.Any("error", err))
		return
	}

	for _, feed := range feeds {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !r.bot.IsModuleEnabled(feed.GuildID, ModuleID) {
			continue
		}

		if err := r.pollSingleFeed(ctx, feed); err != nil {
			r.logger.Warn("failed to poll feed",
				slog.Int64("feed_id", feed.ID),
				slog.String("url", feed.URL),
				slog.Any("error", err),
			)
		}
	}

	// Prune old seen items periodically
	r.pruneOldItems()
}

func (r *RSS) pruneOldItems() {
	cutoff := time.Now().Add(-30 * 24 * time.Hour)
	if err := r.store.PruneSeenItems(cutoff); err != nil {
		r.logger.Warn("failed to prune seen items", slog.Any("error", err))
	}
}
