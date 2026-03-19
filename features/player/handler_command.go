package player

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (p *Player) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ephemeralV2Error("This command can only be used in a server."))
		return
	}

	player := p.lavalink.Player(*guildID)
	queue := p.queues.Get(*guildID)
	ui := BuildPlayerUI(player, queue)

	_ = e.CreateMessage(discord.NewMessageCreateV2(ui))
}
