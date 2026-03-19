package player

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func (p *Player) ensureVoiceConnection(e *events.ModalSubmitInteractionCreate, guildID snowflake.ID) {
	member := e.Member()
	if member == nil {
		return
	}

	voiceState, ok := e.Client().Caches.VoiceState(guildID, member.User.ID)
	if !ok || voiceState.ChannelID == nil {
		return
	}

	_ = e.Client().UpdateVoiceState(context.TODO(), guildID, voiceState.ChannelID, false, true)
}
