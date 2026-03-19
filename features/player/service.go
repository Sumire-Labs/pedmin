package player

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

func (p *Player) handlePause(e *events.ComponentInteractionCreate, guildID snowflake.ID, paused bool) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithPaused(paused)); err != nil {
		p.logger.Error("failed to pause/resume", slog.Any("error", err))
	}
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handleSkip(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}

	queue := p.queues.Get(guildID)
	next, ok := queue.Next()
	if !ok {
		_ = player.Update(context.TODO(), lavalink.WithNullTrack())
		p.respondWithPlayerUpdate(e, player, guildID)
		return
	}

	if err := player.Update(context.TODO(), lavalink.WithTrack(next)); err != nil {
		p.logger.Error("failed to skip", slog.Any("error", err))
	}
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handlePrevious(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}

	queue := p.queues.Get(guildID)
	prev, ok := queue.Previous()
	if !ok {
		if player.Track() != nil {
			_ = player.Update(context.TODO(), lavalink.WithPosition(0))
		}
		p.respondWithPlayerUpdate(e, player, guildID)
		return
	}

	if err := player.Update(context.TODO(), lavalink.WithTrack(prev)); err != nil {
		p.logger.Error("failed to go previous", slog.Any("error", err))
	}
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handleStop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}

	queue := p.queues.Get(guildID)
	queue.Clear()
	_ = player.Update(context.TODO(), lavalink.WithNullTrack())
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handleDisconnect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player != nil {
		_ = player.Destroy(context.TODO())
		p.lavalink.RemovePlayer(guildID)
	}
	p.queues.Delete(guildID)

	_ = e.Client().UpdateVoiceState(context.TODO(), guildID, nil, false, false)

	queue := p.queues.Get(guildID)
	newPlayer := p.lavalink.Player(guildID)
	ui := BuildPlayerUI(newPlayer, queue)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

func (p *Player) handleLoop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	queue.CycleLoop()

	player := p.lavalink.Player(guildID)
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handleVolume(e *events.ComponentInteractionCreate, guildID snowflake.ID, delta int) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}

	newVol := player.Volume() + delta
	if newVol < 0 {
		newVol = 0
	}
	if newVol > 200 {
		newVol = 200
	}

	if err := player.Update(context.TODO(), lavalink.WithVolume(newVol)); err != nil {
		p.logger.Error("failed to set volume", slog.Any("error", err))
	}
	p.respondWithPlayerUpdate(e, player, guildID)
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

	node.LoadTracksHandler(context.TODO(), query, disgolink.NewResultHandler(
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

			p.ensureVoiceConnection(e, guildID)
			player := p.lavalink.Player(guildID)
			_ = player.Update(context.TODO(), lavalink.WithTrack(playlist.Tracks[0]))
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

	if queue.Len() == 1 {
		queue.SetCurrent(0)
	}

	p.ensureVoiceConnection(e, guildID)
	player := p.lavalink.Player(guildID)

	if player.Track() == nil {
		current, ok := queue.Current()
		if !ok {
			current = track
		}
		_ = player.Update(context.TODO(), lavalink.WithTrack(current))
	}
}

func (p *Player) respondWithPlayerUpdate(e *events.ComponentInteractionCreate, player disgolink.Player, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	ui := BuildPlayerUI(player, queue)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

func (p *Player) updatePlayerMessage(player disgolink.Player) {
	queue := p.queues.Get(player.GuildID())
	_ = BuildPlayerUI(player, queue)
}
