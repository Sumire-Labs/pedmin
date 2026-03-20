package ticket

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func (t *Ticket) handleCategorySelect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data := e.Data.(discord.ChannelSelectMenuInteractionData)
	if len(data.Values) == 0 {
		return
	}
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		t.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}
	settings.CategoryID = data.Values[0]
	if err := SaveSettings(t.store, guildID, settings); err != nil {
		t.logger.Error("failed to save settings", slog.Any("error", err))
	}
	_ = e.DeferUpdateMessage()
}

func (t *Ticket) handleLogPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("ログチャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(ModuleID+":log_channel", "ログチャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	).WithEphemeral(true))
}

func (t *Ticket) handleLogChannelSelect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data := e.Data.(discord.ChannelSelectMenuInteractionData)
	if len(data.Values) == 0 {
		return
	}
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		t.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}
	settings.LogChannelID = data.Values[0]
	if err := SaveSettings(t.store, guildID, settings); err != nil {
		t.logger.Error("failed to save settings", slog.Any("error", err))
	}
	_ = e.DeferUpdateMessage()
}

func (t *Ticket) handleRolePrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("サポートロールを選択してください:"),
			discord.NewActionRow(
				discord.NewRoleSelectMenu(ModuleID+":role", "サポートロールを選択..."),
			),
		),
	).WithEphemeral(true))
}

func (t *Ticket) handleRoleSelect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data := e.Data.(discord.RoleSelectMenuInteractionData)
	if len(data.Values) == 0 {
		return
	}
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		t.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}
	settings.SupportRoleID = data.Values[0]
	if err := SaveSettings(t.store, guildID, settings); err != nil {
		t.logger.Error("failed to save settings", slog.Any("error", err))
	}
	_ = e.DeferUpdateMessage()
}
