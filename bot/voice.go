package bot

import (
	"context"

	"github.com/disgoorg/disgo/events"
)

func (b *Bot) onVoiceStateUpdate(e *events.GuildVoiceStateUpdate) {
	if e.VoiceState.UserID != b.Client.ApplicationID {
		return
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), e.VoiceState.GuildID, e.VoiceState.ChannelID, e.VoiceState.SessionID)
}

func (b *Bot) onVoiceServerUpdate(e *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), e.GuildID, e.Token, *e.Endpoint)
}
