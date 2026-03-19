package bot

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
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

func (b *Bot) onMemberVoiceStateUpdate(e *events.GuildVoiceStateUpdate) {
	if e.VoiceState.UserID == b.Client.ApplicationID {
		return
	}
	var channelID snowflake.ID
	if e.VoiceState.ChannelID != nil {
		channelID = *e.VoiceState.ChannelID
	}
	for _, m := range b.Modules {
		if vsl, ok := m.(module.VoiceStateListener); ok {
			vsl.OnVoiceStateUpdate(e.VoiceState.GuildID, channelID, e.VoiceState.UserID)
		}
	}
}
