package bot

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/events"
)

func (b *Bot) onCommandInteraction(e *events.ApplicationCommandInteractionCreate) {
	cmdName := e.SlashCommandInteractionData().CommandName()

	for _, m := range b.Modules {
		for _, cmd := range m.Commands() {
			if cmd.CommandName() == cmdName {
				guildID := e.GuildID()
				if guildID != nil && !b.IsModuleEnabled(*guildID, m.Info().ID) {
					_ = e.CreateMessage(errorMessage("This module is currently disabled."))
					return
				}
				m.HandleCommand(e)
				return
			}
		}
	}
	b.Logger.Warn("unhandled command", slog.String("command", cmdName))
}

func (b *Bot) onComponentInteraction(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	moduleID, _, _ := strings.Cut(customID, ":")

	m, ok := b.Modules[moduleID]
	if !ok {
		b.Logger.Warn("unhandled component", slog.String("custom_id", customID))
		return
	}

	m.HandleComponent(e)
}

func (b *Bot) onModalSubmit(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	moduleID, _, _ := strings.Cut(customID, ":")

	m, ok := b.Modules[moduleID]
	if !ok {
		b.Logger.Warn("unhandled modal", slog.String("custom_id", customID))
		return
	}

	m.HandleModal(e)
}
