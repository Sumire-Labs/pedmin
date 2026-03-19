package player

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (p *Player) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ephemeralV2Error("このコマンドはサーバー内でのみ使用できます。"))
		return
	}

	if err := p.joinVoiceChannel(e.Client(), *guildID, e.Member().User.ID); err != nil {
		_ = e.CreateMessage(ephemeralV2Error("ボイスチャンネルに接続してからコマンドを実行してください。"))
		return
	}

	p.deleteTrackedMessage(*guildID)

	player := p.lavalink.Player(*guildID)
	if player.Volume() != p.defaultVolume {
		ctx, cancel := lavalinkCtx()
		_ = player.Update(ctx, lavalink.WithVolume(p.defaultVolume))
		cancel()
	}
	queue := p.queues.Get(*guildID)
	ui := BuildPlayerUI(player, queue)

	_ = e.CreateMessage(discord.NewMessageCreateV2(ui))

	msg, err := e.Client().Rest.GetInteractionResponse(e.Client().ApplicationID, e.Token())
	if err == nil {
		p.trackMessage(*guildID, msg.ChannelID, msg.ID)
	}
}
