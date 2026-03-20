// Package module defines the Module interface that all feature modules implement.
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

// VoiceStateListener is an optional interface that modules can implement
// to receive voice state updates for non-bot users.
type VoiceStateListener interface {
	OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID)
}
