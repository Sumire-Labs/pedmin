package logger

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type LoggerSettings struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Events    map[string]bool `json:"events"`
}

const (
	EventMessageEdit   = "message_edit"
	EventMessageDelete = "message_delete"
	EventMemberJoin    = "member_join"
	EventMemberLeave   = "member_leave"
	EventBanAdd        = "ban_add"
	EventBanRemove     = "ban_remove"
	EventRoleChange    = "role_change"
	EventChannelChange = "channel_change"
)

var AllEvents = []struct {
	Key   string
	Label string
}{
	{EventMessageEdit, "メッセージ編集"},
	{EventMessageDelete, "メッセージ削除"},
	{EventMemberJoin, "メンバー参加"},
	{EventMemberLeave, "メンバー退出"},
	{EventBanAdd, "BAN"},
	{EventBanRemove, "BAN解除"},
	{EventRoleChange, "ロール変更"},
	{EventChannelChange, "チャンネル変更"},
}

func LoadSettings(guildStore store.GuildStore, guildID snowflake.ID) (*LoggerSettings, error) {
	data, err := guildStore.GetModuleSettings(guildID, ModuleID)
	if err != nil {
		return nil, err
	}
	var s LoggerSettings
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return &LoggerSettings{Events: make(map[string]bool)}, nil
	}
	if s.Events == nil {
		s.Events = make(map[string]bool)
	}
	return &s, nil
}

func SaveSettings(guildStore store.GuildStore, guildID snowflake.ID, settings *LoggerSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return guildStore.SetModuleSettings(guildID, ModuleID, string(data))
}

func (s *LoggerSettings) IsEventEnabled(event string) bool {
	return s.Events[event]
}
