package store

import "github.com/disgoorg/snowflake/v2"

type GuildSettings struct {
	GuildID        snowflake.ID          `json:"guild_id"`
	EnabledModules map[string]bool       `json:"enabled_modules"`
	ModuleSettings map[string]any        `json:"module_settings"`
}

type GuildStore interface {
	Get(guildID snowflake.ID) (*GuildSettings, error)
	Save(settings *GuildSettings) error
	IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
	SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
	Close() error
}
