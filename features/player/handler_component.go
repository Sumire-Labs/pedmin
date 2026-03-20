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

func (p *Player) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	p.logger.Info("component interaction received", slog.String("custom_id", customID))
	_, rest, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		p.logger.Warn("component interaction: guildID is nil")
		return
	}

	action, _, _ := strings.Cut(rest, ":")

	switch action {
	case "skip":
		p.handleSkip(e, *guildID)
	case "stop":
		p.handleStop(e, *guildID)
	case "loop":
		p.handleLoop(e, *guildID)
	case "add":
		p.handleAddModal(e)
	case "queue":
		p.handleShowQueue(e, *guildID)
	case "back":
		p.handleBack(e, *guildID)
	case "clear_queue":
		p.handleClearQueue(e, *guildID)
	}
}

func (p *Player) handleAddModal(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: ModuleID + ":add_modal",
		Title:    "キューに追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("検索キーワードまたはURL",
				discord.NewShortTextInput(ModuleID+":query").
					WithPlaceholder("曲名またはYouTube/SpotifyのURL").
					WithRequired(true),
			),
		},
	})
}

func (p *Player) handleShowQueue(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	player := p.lavalink.Player(guildID)
	ui := BuildQueueUI(queue, player)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

func (p *Player) handleBack(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.Player(guildID)
	queue := p.queues.Get(guildID)
	ui := BuildPlayerUI(player, queue)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

func (p *Player) handleClearQueue(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	queue.Clear()

	player := p.lavalink.Player(guildID)
	ui := BuildQueueUI(queue, player)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
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
		ctx, cancel := lavalinkCtx()
		defer cancel()
		_ = player.Update(ctx, lavalink.WithNullTrack())
		p.respondWithPlayerUpdate(e, player, guildID)
		return
	}

	ctx, cancel := lavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithTrack(next)); err != nil {
		p.logger.Error("failed to skip", slog.Any("error", err))
	}
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handleStop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player != nil {
		ctx, cancel := lavalinkCtx()
		_ = player.Destroy(ctx)
		cancel()
		p.lavalink.RemovePlayer(guildID)
	}
	p.queues.Delete(guildID)
	_ = e.Client().UpdateVoiceState(context.Background(), guildID, nil, false, false)

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

func (p *Player) respondWithPlayerUpdate(e *events.ComponentInteractionCreate, player disgolink.Player, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	ui := BuildPlayerUI(player, queue)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}
