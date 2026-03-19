package player

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// OnVoiceStateUpdate implements module.VoiceStateListener.
func (p *Player) OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID) {
	// Check if the bot is in a VC for this guild
	botVoiceState, ok := p.client.Caches.VoiceState(guildID, p.client.ApplicationID)
	if !ok || botVoiceState.ChannelID == nil {
		return
	}
	botChannelID := *botVoiceState.ChannelID

	// Count non-bot members in the bot's VC
	memberCount := 0
	for vs := range p.client.Caches.VoiceStates(guildID) {
		if vs.ChannelID != nil && *vs.ChannelID == botChannelID && vs.UserID != p.client.ApplicationID {
			memberCount++
		}
	}

	if memberCount == 0 {
		p.startAutoLeaveTimer(guildID)
	} else {
		p.cancelAutoLeaveTimer(guildID)
	}
}

func (p *Player) startAutoLeaveTimer(guildID snowflake.ID) {
	p.cancelAutoLeaveTimer(guildID)

	if p.autoLeaveTimeout == 0 {
		return
	}

	timer := time.AfterFunc(p.autoLeaveTimeout, func() {
		p.logger.Info("auto-leaving voice channel due to inactivity", slog.Any("guild", guildID))
		p.leaveTimers.Delete(guildID)

		// Destroy player
		if player := p.lavalink.ExistingPlayer(guildID); player != nil {
			ctx, cancel := lavalinkCtx()
			_ = player.Destroy(ctx)
			cancel()
			p.lavalink.RemovePlayer(guildID)
		}

		// Disconnect from VC
		_ = p.client.UpdateVoiceState(context.Background(), guildID, nil, false, false)

		// Clear queue
		p.queues.Delete(guildID)

		// Update tracked message to show stopped state
		val, ok := p.messages.Load(guildID)
		if ok {
			tracked := val.(trackedMessage)
			newPlayer := p.lavalink.Player(guildID)
			queue := p.queues.Get(guildID)
			ui := BuildPlayerUI(newPlayer, queue)
			if _, err := p.client.Rest.UpdateMessage(tracked.channelID, tracked.messageID, discord.NewMessageUpdateV2([]discord.LayoutComponent{ui})); err != nil {
				p.logger.Warn("failed to update player message on auto-leave", slog.Any("error", err))
				p.messages.Delete(guildID)
			}
		}
	})

	p.leaveTimers.Store(guildID, timer)
}

func (p *Player) cancelAutoLeaveTimer(guildID snowflake.ID) {
	val, ok := p.leaveTimers.LoadAndDelete(guildID)
	if !ok {
		return
	}
	timer := val.(*time.Timer)
	timer.Stop()
}
