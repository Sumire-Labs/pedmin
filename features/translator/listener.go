package translator

import (
	"context"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func SetupListeners(client *disgobot.Client, t *Translator) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(t.onMessageReactionAdd),
	)
}

func (t *Translator) onMessageReactionAdd(e *events.GuildMessageReactionAdd) {
	// Ignore bot reactions
	if e.Member.User.Bot {
		return
	}

	if !t.bot.IsModuleEnabled(e.GuildID, ModuleID) {
		return
	}

	if !t.translateClient.IsAvailable() {
		return
	}

	// Check if emoji is a flag
	if e.Emoji.Name == nil {
		return
	}
	targetLang, ok := flagToLang[*e.Emoji.Name]
	if !ok {
		return
	}

	t.processTranslation(context.Background(), e.ChannelID, e.MessageID, targetLang)
}
