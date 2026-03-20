package store

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type GuildSettings struct {
	GuildID        snowflake.ID          `json:"guild_id"`
	EnabledModules map[string]bool       `json:"enabled_modules"`
	ModuleSettings map[string]any        `json:"module_settings"`
}

type Ticket struct {
	GuildID   snowflake.ID
	Number    int
	ChannelID snowflake.ID
	UserID    snowflake.ID
	Subject   string
	CreatedAt time.Time
	ClosedAt  *time.Time
	ClosedBy  *snowflake.ID
}

type GuildStore interface {
	Get(guildID snowflake.ID) (*GuildSettings, error)
	Save(settings *GuildSettings) error
	IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
	SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
	GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error)
	SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error
	CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error
	GetTicketByChannel(channelID snowflake.ID) (*Ticket, error)
	CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error
	DeleteTicket(channelID snowflake.ID) error
	Close() error
}
