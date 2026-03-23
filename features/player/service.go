package player

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

func (p *Player) lavalinkCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), p.lavalinkTimeout)
}

func (p *Player) loadAndPlay(e *events.ModalSubmitInteractionCreate, guildID snowflake.ID, query string) {
	if !strings.HasPrefix(query, "http") {
		query = "ytsearch:" + query
	}

	node := p.lavalink.BestNode()
	if node == nil {
		p.logger.Error("no lavalink node available")
		p.sendLoadFollowup(e, "❌ 音楽サーバーに接続できません。")
		return
	}

	loadCtx, loadCancel := context.WithTimeout(context.Background(), p.lavalinkLoadTimeout)
	defer loadCancel()
	node.LoadTracksHandler(loadCtx, query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			p.playTrack(e, guildID, track)
			p.sendLoadFollowup(e, fmt.Sprintf("🎵 キューに追加: **%s**", track.Info.Title))
		},
		func(playlist lavalink.Playlist) {
			if len(playlist.Tracks) == 0 {
				p.sendLoadFollowup(e, "❌ プレイリストにトラックがありません。")
				return
			}
			queue := p.queues.Get(guildID)
			queue.Add(playlist.Tracks...)
			queue.SetCurrent(0)

			_ = p.joinVoiceChannel(e.Client(), guildID, e.Member().User.ID)
			player := p.lavalink.Player(guildID)
			ctx, cancel := p.lavalinkCtx()
			_ = player.Update(ctx, lavalink.WithTrack(playlist.Tracks[0]))
			cancel()
			p.sendLoadFollowup(e, fmt.Sprintf("🎵 プレイリストから %d 曲を追加しました", len(playlist.Tracks)))
		},
		func(tracks []lavalink.Track) {
			if len(tracks) == 0 {
				p.sendLoadFollowup(e, "❌ 検索結果が見つかりません。")
				return
			}
			p.playTrack(e, guildID, tracks[0])
			p.sendLoadFollowup(e, fmt.Sprintf("🎵 キューに追加: **%s**", tracks[0].Info.Title))
		},
		func() {
			p.logger.Info("no matches found", slog.String("query", query))
			p.sendLoadFollowup(e, "❌ 検索結果が見つかりません。")
		},
		func(err error) {
			p.logger.Error("load failed", slog.Any("error", err))
			p.sendLoadFollowup(e, "❌ トラックの読み込みに失敗しました。")
		},
	))
}

func (p *Player) sendLoadFollowup(e *events.ModalSubmitInteractionCreate, text string) {
	_, _ = e.Client().Rest.UpdateInteractionResponse(
		e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(discord.NewTextDisplay(text)),
		}),
	)
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
		ctx, cancel := p.lavalinkCtx()
		_ = player.Update(ctx, lavalink.WithTrack(current))
		cancel()
	}
}
