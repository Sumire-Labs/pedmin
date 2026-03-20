package ticket

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type TicketSettings struct {
	CategoryID    snowflake.ID `json:"category_id"`
	LogChannelID  snowflake.ID `json:"log_channel_id"`
	SupportRoleID snowflake.ID `json:"support_role_id"`
	NextNumber    int          `json:"next_number"`
}

func LoadSettings(guildStore store.GuildStore, guildID snowflake.ID) (*TicketSettings, error) {
	data, err := guildStore.GetModuleSettings(guildID, ModuleID)
	if err != nil {
		return nil, err
	}
	var s TicketSettings
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return &TicketSettings{}, nil
	}
	return &s, nil
}

func SaveSettings(guildStore store.GuildStore, guildID snowflake.ID, s *TicketSettings) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return guildStore.SetModuleSettings(guildID, ModuleID, string(data))
}
