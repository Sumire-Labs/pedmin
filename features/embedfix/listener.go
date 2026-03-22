package embedfix

import (
	"context"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func SetupListeners(client *disgobot.Client, ef *EmbedFix) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(ef.onMessageCreate),
	)
}

func (ef *EmbedFix) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if !ef.bot.IsModuleEnabled(e.GuildID, ModuleID) {
		return
	}

	ef.processMessageURLs(context.Background(), e)
}
