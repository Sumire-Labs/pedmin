package player

import (
	"context"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

const progressTickInterval = 10 * time.Second

// startProgressTicker starts a per-guild ticker that periodically updates
// the player message to refresh the progress bar. Calling it when a ticker
// already exists for the guild cancels the old one first.
func (p *Player) startProgressTicker(guildID snowflake.ID) {
	p.stopProgressTicker(guildID)

	ctx, cancel := context.WithCancel(context.Background())
	p.progressTickers.Store(guildID, cancel)

	go p.progressTickLoop(ctx, guildID)
}

func (p *Player) stopProgressTicker(guildID snowflake.ID) {
	val, ok := p.progressTickers.LoadAndDelete(guildID)
	if !ok {
		return
	}
	cancel, ok := val.(context.CancelFunc)
	if !ok {
		return
	}
	cancel()
}

// stopAllProgressTickers cancels all active tickers. Called during shutdown.
func (p *Player) stopAllProgressTickers() {
	p.progressTickers.Range(func(key, value any) bool {
		if cancel, ok := value.(context.CancelFunc); ok {
			cancel()
		}
		p.progressTickers.Delete(key)
		return true
	})
}

func (p *Player) progressTickLoop(ctx context.Context, guildID snowflake.ID) {
	ticker := time.NewTicker(progressTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			player := p.lavalink.ExistingPlayer(guildID)
			if player == nil || player.Track() == nil || player.Paused() {
				continue
			}
			p.updatePlayerMessage(player)
		}
	}
}
