package player

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

const lavalinkTimeout = 2 * time.Second

func lavalinkCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), lavalinkTimeout)
}

func (p *Player) loadAndPlay(e *events.ModalSubmitInteractionCreate, guildID snowflake.ID, query string) {
	if !strings.HasPrefix(query, "http") {
		query = "ytsearch:" + query
	}

	node := p.lavalink.BestNode()
	if node == nil {
		p.logger.Error("no lavalink node available")
		return
	}

	loadCtx, loadCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer loadCancel()
	node.LoadTracksHandler(loadCtx, query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			p.playTrack(e, guildID, track)
		},
		func(playlist lavalink.Playlist) {
			if len(playlist.Tracks) == 0 {
				return
			}
			queue := p.queues.Get(guildID)
			queue.Add(playlist.Tracks...)
			queue.SetCurrent(0)

			_ = p.joinVoiceChannel(e.Client(), guildID, e.Member().User.ID)
			player := p.lavalink.Player(guildID)
			ctx, cancel := lavalinkCtx()
			_ = player.Update(ctx, lavalink.WithTrack(playlist.Tracks[0]))
			cancel()
		},
		func(tracks []lavalink.Track) {
			if len(tracks) == 0 {
				return
			}
			p.playTrack(e, guildID, tracks[0])
		},
		func() {
			p.logger.Info("no matches found", slog.String("query", query))
		},
		func(err error) {
			p.logger.Error("load failed", slog.Any("error", err))
		},
	))
}

func (p *Player) playTrack(e *events.ModalSubmitInteractionCreate, guildID snowflake.ID, track lavalink.Track) {
	queue := p.queues.Get(guildID)
	queue.Add(track)

	// Only set current when this is the first track; otherwise the queue
	// already has a current track and this one should wait its turn.
	if queue.Len() == 1 {
		queue.SetCurrent(0)
	}

	_ = p.joinVoiceChannel(e.Client(), guildID, e.Member().User.ID)
	player := p.lavalink.Player(guildID)

	if player.Track() == nil {
		current, ok := queue.Current()
		if !ok {
			current = track
		}
		ctx, cancel := lavalinkCtx()
		_ = player.Update(ctx, lavalink.WithTrack(current))
		cancel()
	}
}
