package player

import (
	"github.com/disgoorg/disgo/events"
)

func (p *Player) HandleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	if customID != ModuleID+":add_modal" {
		return
	}

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	ti, ok := e.Data.TextInput(ModuleID + ":query")
	query := ""
	if ok {
		query = ti.Value
	}

	if query == "" {
		_ = e.CreateMessage(ephemeralV2Error("Please provide a search query or URL."))
		return
	}

	_ = e.DeferCreateMessage(true)
	p.loadAndPlay(e, *guildID, query)
}
