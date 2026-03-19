package module

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type Info struct {
	ID          string
	Name        string
	Description string
	AlwaysOn    bool
}

type Module interface {
	Info() Info
	Commands() []discord.ApplicationCommandCreate
	HandleCommand(e *events.ApplicationCommandInteractionCreate)
	HandleComponent(e *events.ComponentInteractionCreate)
	HandleModal(e *events.ModalSubmitInteractionCreate)
	SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent
	HandleSettingsComponent(e *events.ComponentInteractionCreate)
}
